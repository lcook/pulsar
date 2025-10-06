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

	session.ChannelMessageSendEmbed(message.BeforeDelete.ChannelID, &discordgo.MessageEmbed{
		Color:       embedDeleteColor,
		Description: fmt.Sprintf("Message sent by <@!%s> in <#%s> deleted", message.BeforeDelete.Author.ID, message.BeforeDelete.ChannelID),
		Author: &discordgo.MessageEmbedAuthor{
			Name:    message.BeforeDelete.Author.Username,
			IconURL: message.BeforeDelete.Author.AvatarURL("96"),
		},
		Fields:    []*discordgo.MessageEmbedField{{Value: buildContentField(message.BeforeDelete.Content, message.BeforeDelete.Attachments)}},
		Timestamp: message.BeforeDelete.Timestamp.Format(time.RFC3339),
		Footer:    &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s", message.BeforeDelete.ID)},
	})
}
