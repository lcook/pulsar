// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
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

	logUser(m.User, log.DebugLevel, "GuildMemberRemove(event): Member left")
}
