/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package event

import (
	"crypto/sha512"
	"encoding/hex"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Member == nil || m.Author.Bot {
		return
	}

	if !h.Settings.AntiSpamSettings.Enabled {
		return
	}

	for _, id := range h.Settings.AntiSpamSettings.ExcludedRoleIDs {
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

	sha := sha512.New()
	sha.Write([]byte(content.String()))

	hash := hex.EncodeToString(sha.Sum(nil))

	h.Logs.Add(Log{Message: m.Message, Hash: hash})

	if len(h.Logs.Slice()) > 1 {
		h.AntiSpam(s, m, hash)
	}
}
