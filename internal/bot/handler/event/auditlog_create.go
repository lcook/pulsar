// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) AuditLogCreate(
	s *discordgo.Session,
	e *discordgo.GuildAuditLogEntryCreate,
) {
	var (
		icon     = ":hammer:"
		verb     string
		duration time.Duration
	)

	switch *e.ActionType {
	case discordgo.AuditLogActionMemberKick:
		verb = "kicked"
	case discordgo.AuditLogActionMemberBanAdd:
		verb = "banned"
	case discordgo.AuditLogActionMemberUpdate:
		var found bool

		for _, val := range e.Changes {
			if *val.Key != discordgo.AuditLogChangeKeyCommunicationDisabledUntil ||
				val.NewValue == nil {
				continue
			}

			value, err := time.Parse(time.RFC3339, val.NewValue.(string))
			if err != nil {
				continue
			}

			timestamp, err := discordgo.SnowflakeTimestamp(e.ID)
			if err != nil {
				continue
			}

			duration = value.Sub(timestamp).Abs()
			found = true

			break
		}

		if !found {
			return
		}

		icon = ":mute:"
		verb = "timed out"
	default:
		return
	}

	var (
		fields    = make([]*discordgo.MessageEmbedField, 0, 2)
		logFields = make([]log.Fields, 0, 2)
	)

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

	if verb == "timed out" && duration > 0 {
		logFields = append(
			logFields,
			log.Fields{"duration": duration.String()},
		)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Duration",
			Value:  duration.String(),
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
				"%s **User %s has been %s**",
				icon,
				user.Mention(),
				verb,
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
		fmt.Sprintf("AuditLogCreate(event): User %s", verb),
		logFields...,
	)

	h.ForwardAlert(s, message, false)
}
