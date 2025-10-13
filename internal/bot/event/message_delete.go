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
)

func MessageDelete(session *discordgo.Session, message *discordgo.MessageDelete) {
	if message.BeforeDelete == nil ||
		message.BeforeDelete.Author.ID == session.State.User.ID {
		return
	}

	session.ChannelMessageSendEmbed(eventLogChannel, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("**:wastebasket: Message sent by <@!%s> in <#%s> deleted**", message.BeforeDelete.Author.ID, message.BeforeDelete.ChannelID),
		Timestamp:   message.BeforeDelete.Timestamp.Format(time.RFC3339),
		Color:       embedDeleteColor,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s", message.BeforeDelete.ID)},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    message.BeforeDelete.Author.Username,
			IconURL: message.BeforeDelete.Author.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{{
			Name:  "Contents",
			Value: buildContentField(message.BeforeDelete.Content, message.BeforeDelete.Attachments),
		}},
	})
}
