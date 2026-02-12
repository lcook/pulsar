// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) MessageUpdate(
	s *discordgo.Session,
	m *discordgo.MessageUpdate,
) {
	if m.BeforeUpdate == nil {
		return
	}

	if m.Author.ID == s.State.User.ID || m.Member == nil || m.Author.Bot {
		return
	}

	if m.BeforeUpdate.Content == m.Content &&
		len(m.BeforeUpdate.Attachments) == len(m.Attachments) {
		return
	}

	link := fmt.Sprintf(
		"%schannels/%s/%s/%s",
		discordgo.EndpointDiscord,
		m.GuildID,
		m.ChannelID,
		m.ID,
	)

	if canViewChannel(s, m.GuildID, m.ChannelID) {
		s.ChannelMessageSendEmbed(
			h.Settings.LogChannel,
			&discordgo.MessageEmbed{
				Description: fmt.Sprintf(
					"**:pencil: [Message](%s) edited by %s in <#%s>**",
					link,
					m.Author.Mention(),
					m.ChannelID,
				),
				Color: embedUpdateColor,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.Author.Username,
					IconURL: m.Author.AvatarURL("256"),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name: "Before",
						Value: buildContentField(
							m.BeforeUpdate.Content,
							m.BeforeUpdate.Attachments,
							m.BeforeUpdate.StickerItems,
						),
						Inline: true,
					},
					{
						Name: "After",
						Value: buildContentField(
							m.Content,
							m.Attachments,
							m.StickerItems,
						),
						Inline: true,
					},
				},
			},
		)
	}
}
