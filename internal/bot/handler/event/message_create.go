// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/lcook/pulsar/internal/antispam"
)

func (h *Handler) MessageCreate(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Author.ID == s.State.User.ID || m.Member == nil || m.Author.Bot {
		return
	}

	if !h.Settings.Enabled {
		return
	}

	for _, id := range h.Settings.ExcludedRoleIDs {
		if slices.Contains(m.Member.Roles, id) {
			return
		}
	}

	var content strings.Builder
	for _, attachment := range m.Attachments {
		content.WriteString(attachment.Filename)
	}

	if m.Content != "" {
		content.WriteString(m.Content)
	}

	if content.String() == "" {
		return
	}

	hash := hashContent(content.String())

	h.Logs.Add(antispam.Log{Message: m.Message, Hash: hash})

	if len(h.Logs.Slice()) > 1 {
		logs, rule := antispam.Run(m, hash, h.Logs, h.Settings.Rules)
		if len(logs) == 0 || rule == nil {
			return
		}

		h.ProcessSpam(s, m, logs, rule)
	}
}
