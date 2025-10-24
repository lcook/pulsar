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

func (h *Handler) AutoModExecution(s *discordgo.Session, am *discordgo.AutoModerationActionExecution) {
	if am.Action.Type != discordgo.AutoModerationRuleActionBlockMessage {
		return
	}

	user, _ := s.User(am.UserID)

	s.ChannelMessageSendEmbed(h.Settings.LogChannel, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("**:tools: AutoMod action triggered**: message sent by <@!%s> in <#%s> flagged. _Please do not click any links it may contain as they may be dangerous_", user.ID, am.ChannelID),
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       embedDeleteColor,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Rule: %s", am.RuleID)},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.Username,
			IconURL: user.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Contents",
				Value: truncateContent(am.Content),
			},
		},
	})
}
