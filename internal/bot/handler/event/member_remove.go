/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package event

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) GuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if m.User == nil {
		return
	}

	/*
	 * Discord's audit log entries sometimes appear *after* the member remove event.
	 * Introduce a small delay to try and circumvent this.
	 */
	time.Sleep(1 * time.Second)

	entries, err := auditLogActions(s, m.Member, 0, 15)
	if err != nil || len(entries) < 1 {
		return
	}

	fields := make([]*discordgo.MessageEmbedField, 0, 2)
	now := time.Now().UTC()

	var action string

	for _, entry := range entries {
		switch *entry.ActionType {
		case discordgo.AuditLogActionMemberKick:
			action = "kicked"
		case discordgo.AuditLogActionMemberBanAdd:
			action = "banned"
		default:
			continue
		}

		timestamp, _ := discordgo.SnowflakeTimestamp(entry.ID)
		/*
		 * Filter to entries within the last 30 seconds to prevent false positives.
		 *
		 * Without this check, when a user:
		 *   1. Is kicked/banned (audit entry created)
		 *   2. Rejoins the server
		 *   3. Voluntarily leaves later
		 *
		 * The bot would incorrectly report them as kicked/banned using the old audit entry.
		 *
		 * This ensures we only consider actions that could logically have caused
		 * the current removal event.
		 */
		if now.Sub(timestamp.UTC()) > 60*time.Second {
			continue
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Moderator",
			Value:  fmt.Sprintf("<@!%s>", entry.UserID),
			Inline: true,
		})

		if entry.Reason != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Reason",
				Value:  entry.Reason,
				Inline: true,
			})
		}

		break
	}

	if action != "" {
		s.ChannelMessageSendEmbed(h.Settings.LogChannel, &discordgo.MessageEmbed{
			Description: fmt.Sprintf(":hammer: **Member <@!%s> has been %s**", m.User.ID, action),
			Timestamp:   time.Now().Format(time.RFC3339),
			Color:       embedDeleteColor,
			Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s", m.User.ID)},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    m.User.Username,
				IconURL: m.AvatarURL("256"),
			},
			Fields: fields,
		})
	}
}
