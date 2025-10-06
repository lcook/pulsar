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

	session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
		Color:       embedUpdateColor,
		Description: fmt.Sprintf("Message sent by <@!%s> in <#%s> updated", message.Author.ID, message.ChannelID),
		Author: &discordgo.MessageEmbedAuthor{
			Name:    message.Author.Username,
			IconURL: message.Author.AvatarURL("96"),
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
		Timestamp: message.EditedTimestamp.Format(time.RFC3339),
		Footer:    &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s", message.ID)},
	})
}
