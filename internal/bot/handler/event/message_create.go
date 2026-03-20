// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"

	"github.com/lcook/pulsar/internal/antispam"
)

func (h *Handler) MessageCreate(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Author.ID == s.State.User.ID || m.Member == nil || m.Author.Bot {
		return
	}

	if m.Type != discordgo.MessageTypeDefault {
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
	for _, sticker := range m.StickerItems {
		content.WriteString(sticker.Name)
	}

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
	log.WithFields(log.Fields{
		"author_id":    m.Author.ID,
		"channel_id":   m.ChannelID,
		"content_hash": hash[0:12],
	}).Trace("MessageCreate: content hashed for antispam analysis")

	if len(h.Logs.Slice()) > 1 {
		logs, rule := antispam.Run(m, hash, h.Logs, h.Settings.Rules)
		if len(logs) == 0 || rule == nil {
			return
		}

		h.ProcessSpam(s, m, logs, rule)
	}
}
