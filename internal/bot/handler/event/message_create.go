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
	if m.Author.ID == s.State.User.ID || m.Member == nil || m.Author.Bot || m.Content == "" {
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

/*
	panic: runtime error: invalid memory address or nil pointer dereference
	[signal SIGSEGV: segmentation violation code=0x1 addr=0x68 pc=0x7dc447]

	goroutine 6316 [running]:
	github.com/lcook/pulsar/internal/bot/handler/event.(*Handler).MessageCreate(0x870126400, 0x87011ae08, 0x8702065b8)
	        /home/lcook/dev/pulsar/internal/bot/handler/event/message_create.go:26 +0x107
	github.com/bwmarrin/discordgo.messageCreateEventHandler.Handle(0x87003e008?, 0x8706b47d0?, {0x940980?, 0x8702065b8?})
	        /home/lcook/go/pkg/mod/github.com/bwmarrin/discordgo@v0.29.0/eventhandlers.go:926 +0x33
	created by github.com/bwmarrin/discordgo.(*Session).handle in goroutine 5935
	        /home/lcook/go/pkg/mod/github.com/bwmarrin/discordgo@v0.29.0/event.go:171 +0x12c

*/
