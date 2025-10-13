package event

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func AutoModExecution(session *discordgo.Session, automod *discordgo.AutoModerationActionExecution) {
	if automod.Action.Type != discordgo.AutoModerationRuleActionBlockMessage {
		return
	}

	user, _ := session.User(automod.UserID)

	session.ChannelMessageSendEmbed(eventLogChannel, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("**:tools: AutoMod action triggered**: message sent by <@!%s> in <#%s> flagged. _Please do not click any links it may contain as they may be dangerous_", user.ID, automod.ChannelID),
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       embedDeleteColor,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Rule: %s", automod.RuleID)},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.Username,
			IconURL: user.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Contents",
				Value: truncateContent(automod.Content),
			},
		},
	})
}
