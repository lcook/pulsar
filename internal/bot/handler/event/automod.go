// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) AutoModExecution(
	s *discordgo.Session,
	am *discordgo.AutoModerationActionExecution,
) {
	if am.Action.Type != discordgo.AutoModerationRuleActionBlockMessage {
		return
	}

	user, err := s.User(am.UserID)
	if err != nil {
		return
	}

	message, _ := s.ChannelMessageSendEmbed(
		h.Settings.LogChannel,
		&discordgo.MessageEmbed{
			Title: fmt.Sprintf(":shield: AutoMod alert (%s)", am.RuleID),
			Description: fmt.Sprintf(
				"Message sent by %s in <#%s> flagged by AutoMod. _Please exercise caution: these messages may contain malicious links, phishing attempts, or other harmful content_.",
				user.Mention(),
				am.ChannelID,
			),
			Color: embedDeleteColor,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    user.Username,
				IconURL: user.AvatarURL("256"),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Content",
					Value: TruncateContent(am.Content),
				},
			},
		},
	)

	h.ForwardAlert(s, message, false)
}
