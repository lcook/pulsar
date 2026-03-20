// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
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
		h.Errors <- HandlerChannel{
			Message: "AutoModExecution(event): Unable to fetch user information",
			Fields: log.Fields{
				"user_id":       am.UserID,
				"error_message": err.Error(),
			},
		}

		return
	}

	message, err := sendSilentEmbed(s, h.Settings.LogChannel,
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
	if err != nil {
		h.Errors <- HandlerChannel{
			Message: "AutoModExecution(event): Unable to send message embed",
			Fields: log.Fields{
				"message_id":    message.ID,
				"error_message": err.Error(),
			},
		}

		return
	}

	h.ForwardAlert(s, message, false)
}
