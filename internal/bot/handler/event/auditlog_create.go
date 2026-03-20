// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) AuditLogCreate(
	s *discordgo.Session,
	e *discordgo.GuildAuditLogEntryCreate,
) {
	actionType := *e.ActionType
	if actionType != discordgo.AuditLogActionMemberBanAdd &&
		actionType != discordgo.AuditLogActionMemberKick {
		return
	}

	action := "banned"
	if actionType == discordgo.AuditLogActionMemberKick {
		action = "kicked"
	}

	fields := make([]*discordgo.MessageEmbedField, 0, 2)
	logFields := make([]log.Fields, 0, 2)

	logFields = append(logFields, log.Fields{"moderator": e.UserID})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Moderator",
		Value:  fmt.Sprintf("<@!%s>", e.UserID),
		Inline: true,
	})

	if e.Reason != "" {
		logFields = append(
			logFields,
			log.Fields{"reason": e.Reason},
		)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Reason",
			Value:  e.Reason,
			Inline: true,
		})
	}

	user, err := s.User(e.TargetID)
	if err != nil {
		h.Errors <- HandlerChannel{
			Message: "AuditLogCreate(event): Unable to fetch user information",
			Fields: log.Fields{
				"user_id":       e.TargetID,
				"error_message": err.Error(),
			},
		}

		return
	}

	message, err := sendSilentEmbed(
		s,
		h.Settings.LogChannel,
		&discordgo.MessageEmbed{
			Description: fmt.Sprintf(
				":hammer: **User %s has been %s**",
				user.Mention(),
				action,
			),
			Color: embedDeleteColor,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    user.Username,
				IconURL: user.AvatarURL("256"),
			},
			Fields: fields,
		},
	)
	if err != nil {
		h.Errors <- HandlerChannel{
			Message: "AuditLogCreate(event): Unable to send message embed",
			Fields: log.Fields{
				"message_id":    message.ID,
				"error_message": err.Error(),
			},
		}

		return
	}

	logUser(
		user,
		log.WarnLevel,
		fmt.Sprintf("AuditLogCreate(event): User %s", action),
		logFields...,
	)

	h.ForwardAlert(s, message, false)
}
