/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != h.Settings.Prefix+"help" {
		return
	}

	fields := make([]*discordgo.MessageEmbedField, 0, len(h.commands))
	for _, command := range h.commands {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  h.Settings.Prefix + command.Name,
			Value: command.Description,
		})
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("List of commands available (%d) to use in the server.", len(h.commands)),
		Color:       embedColorFreeBSD,
		Fields:      fields,
	})
}
