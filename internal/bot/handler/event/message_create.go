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

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot || m.Content == "" {
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

	sha := sha512.New()
	sha.Write([]byte(m.Content))

	hash := hex.EncodeToString(sha.Sum(nil))

	h.Logs.Add(Log{Message: m.Message, Hash: hash})

	if len(h.Logs.Slice()) > 1 {
		h.AntiSpam(s, m, hash)
	}
}
