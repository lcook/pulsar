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

func (h *Handler) GuildMemberAdd(
	s *discordgo.Session,
	m *discordgo.GuildMemberAdd,
) {
	logUser(m.User, log.DebugLevel, "Member joined")

	created, _ := discordgo.SnowflakeTimestamp(m.User.ID)

	age := m.JoinedAt.UTC().Sub(created.UTC())

	if age <= h.Settings.MinumumAccountAge {
		logUser(
			m.User,
			log.WarnLevel,
			"Suspected spam or advertising account joined",
		)

		message, _ := s.ChannelMessageSendEmbed(
			h.Settings.LogChannel,
			&discordgo.MessageEmbed{
				Title: ":shield: Suspected spam or advertising account",
				Description: fmt.Sprintf(
					"User %s joined with a recently created account, it may be used for spam or advertising - exercise caution.",
					m.Mention(),
				),
				Timestamp: time.Now().Format(time.RFC3339),
				Color:     embedUpdateColor,
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("ID: %s", m.User.ID),
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.User.Username,
					IconURL: m.User.AvatarURL("256"),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Account age",
						Value: age.String(),
					},
				},
			},
		)

		h.ForwardAlert(s, message, false)
	}
}
