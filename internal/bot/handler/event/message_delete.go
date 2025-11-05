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
	"github.com/lcook/pulsar/internal/antispam"
)

func (h *Handler) MessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.BeforeDelete == nil {
		return
	}

	if m.BeforeDelete.Author.ID == s.State.User.ID {
		return
	}

	var spam bool

	h.Logs.ForEach(func(l *antispam.Log) {
		if l.Message.ID == m.ID && l.Deleted() {
			spam = true
		}
	})

	if !spam && canViewChannel(s, m.GuildID, m.ChannelID) {
		s.ChannelMessageSendEmbed(h.Settings.LogChannel, &discordgo.MessageEmbed{
			Description: fmt.Sprintf("**:wastebasket: Message deleted by <@!%s> in <#%s>**", m.BeforeDelete.Author.ID, m.BeforeDelete.ChannelID),
			Timestamp:   m.BeforeDelete.Timestamp.Format(time.RFC3339),
			Color:       embedDeleteColor,
			Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s", m.BeforeDelete.ID)},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    m.BeforeDelete.Author.Username,
				IconURL: m.BeforeDelete.Author.AvatarURL("256"),
			},
			Fields: []*discordgo.MessageEmbedField{{
				Name:  "Contents",
				Value: buildContentField(m.BeforeDelete.Content, m.BeforeDelete.Attachments),
			}},
		})
	}
}
