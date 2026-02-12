// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"fmt"

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

		message, _ := sendSilentEmbed(
			s,
			h.Settings.LogChannel,
			&discordgo.MessageEmbed{
				Title: ":rotating_light: New account detected",
				Description: fmt.Sprintf(
					"-# Attention: User %s has joined with a newly created account (%s). This might be a spam, advertising or compromised account - exercise caution.",
					m.Mention(),
					age.String(),
				),
				Color: embedUpdateColor,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.User.Username,
					IconURL: m.User.AvatarURL("256"),
				},
			},
		)

		h.ForwardAlert(s, message, false)
	}
}
