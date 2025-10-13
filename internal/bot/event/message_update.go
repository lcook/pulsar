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

func MessageUpdate(session *discordgo.Session, message *discordgo.MessageUpdate) {
	if message.Author.ID == session.State.User.ID ||
		message.BeforeUpdate == nil {
		return
	}

	if message.BeforeUpdate.Content == message.Content &&
		len(message.BeforeUpdate.Attachments) == len(message.Attachments) {
		return
	}

	link := fmt.Sprintf("%schannels/%s/%s/%s", discordgo.EndpointDiscord, message.GuildID, message.ChannelID, message.ID)

	session.ChannelMessageSendEmbed(eventLogChannel, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("**:pencil: [Message](%s) edited by <@!%s> in <#%s>**", link, message.Author.ID, message.ChannelID),
		Timestamp:   message.EditedTimestamp.Format(time.RFC3339),
		Color:       embedUpdateColor,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s", message.ID)},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    message.Author.Username,
			IconURL: message.Author.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Before",
				Value:  buildContentField(message.BeforeUpdate.Content, message.BeforeUpdate.Attachments),
				Inline: true,
			},
			{
				Name:   "After",
				Value:  buildContentField(message.Content, message.Attachments),
				Inline: true,
			},
		},
	})
}
