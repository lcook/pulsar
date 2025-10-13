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

func GuildMemberRemove(session *discordgo.Session, member *discordgo.GuildMemberRemove) {
	if member.User == nil {
		return
	}

	entries, err := auditLogActions(session, member.Member, 0)
	if err != nil {
		return
	}

	fields := make([]*discordgo.MessageEmbedField, 0, 2)

	var action string

	for _, entry := range entries {
		switch *entry.ActionType {
		case discordgo.AuditLogActionMemberKick:
			action = "kicked"
		case discordgo.AuditLogActionMemberBanAdd:
			action = "banned"
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

	session.ChannelMessageSendEmbed(eventLogChannel, &discordgo.MessageEmbed{
		Description: fmt.Sprintf(":hammer: **Member <@!%s> has been %s**", member.User.ID, action),
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       embedDeleteColor,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s", member.User.ID)},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    member.User.Username,
			IconURL: member.AvatarURL("256"),
		},
		Fields: fields,
	})
}
