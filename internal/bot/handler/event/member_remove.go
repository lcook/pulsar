// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) GuildMemberRemove(
	s *discordgo.Session,
	m *discordgo.GuildMemberRemove,
) {
	if m.User == nil {
		return
	}

	logMember(m.User, log.TraceLevel, "Member left")

	entries, err := auditLogActionsLast(s, m.Member, 0, 15, 60*time.Second)
	if err != nil || len(entries) < 1 {
		return
	}

	var action string

	fields := make([]*discordgo.MessageEmbedField, 0, 2)
	logFields := make([]log.Fields, 0, 2)

	for _, entry := range entries {
		switch *entry.ActionType {
		case discordgo.AuditLogActionMemberKick:
			action = "kicked"
		case discordgo.AuditLogActionMemberBanAdd:
			action = "banned"
		default:
			continue
		}

		logFields = append(logFields, log.Fields{"moderator": entry.UserID})

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Moderator",
			Value:  fmt.Sprintf("<@!%s>", entry.UserID),
			Inline: true,
		})

		if entry.Reason != "" {
			logFields = append(
				logFields,
				log.Fields{"reason": entry.Reason},
			)

			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Reason",
				Value:  entry.Reason,
				Inline: true,
			})
		}

		break
	}

	if action != "" {
		message, _ := s.ChannelMessageSendEmbed(
			h.Settings.LogChannel,
			&discordgo.MessageEmbed{
				Description: fmt.Sprintf(
					":hammer: **Member <@!%s> has been %s**",
					m.User.ID,
					action,
				),
				Timestamp: time.Now().Format(time.RFC3339),
				Color:     embedDeleteColor,
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("ID: %s", m.User.ID),
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.User.Username,
					IconURL: m.AvatarURL("256"),
				},
				Fields: fields,
			},
		)

		logMember(
			m.User,
			log.WarnLevel,
			fmt.Sprintf("Member %s", action),
			logFields...,
		)

		h.ForwardAlert(s, message, false)
	}
}
