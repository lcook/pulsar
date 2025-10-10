/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package command

import (
	"github.com/bwmarrin/discordgo"
)

var Help = Command{
	"help",
	"Show help page",
	func(session *discordgo.Session, message *discordgo.MessageCreate) {
		if message.Author.ID == session.State.User.ID || message.Content != "!help" {
			return
		}

		var fields []*discordgo.MessageEmbedField
		for _, command := range []Command{Role, Bug} {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  prefix + command.Name,
				Value: command.Description,
			})
		}

		session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
			Description: "List of commands available to use in the server",
			Color:       embedColorFreeBSD,
			Fields:      fields,
		})
	},
}
