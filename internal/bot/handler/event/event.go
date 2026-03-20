// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"

	"github.com/lcook/pulsar/internal/antispam"
	"github.com/lcook/pulsar/internal/cache"
	"github.com/lcook/pulsar/internal/config"
)

const (
	embedDeleteColor int = 0xDC322F
	embedUpdateColor int = 0x268BD2
)

type HandlerChannel struct {
	Message string
	Fields  log.Fields
}

type Handler struct {
	Settings config.Settings
	Events   []any
	Logs     *cache.RingBuffer[antispam.Log]
	Errors   chan HandlerChannel
}

func New(settings config.Settings, buffer uint64) *Handler {
	h := &Handler{
		Settings: settings,
		Logs:     cache.NewRingBuffer[antispam.Log](buffer),
		Errors:   make(chan HandlerChannel),
	}

	h.Events = append(h.Events, h.MessageCreate)
	h.Events = append(h.Events, h.MessageDelete)
	h.Events = append(h.Events, h.MessageUpdate)
	h.Events = append(h.Events, h.GuildMemberAdd)
	h.Events = append(h.Events, h.GuildMemberRemove)
	h.Events = append(h.Events, h.AutoModExecution)
	h.Events = append(h.Events, h.AuditLogCreate)

	return h
}

func (h *Handler) ProcessSpam(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
	logs []*antispam.Log,
	rule *antispam.HeuristicRule,
) {
	timeout := time.Now().Add(rule.Timeout)

	if err := session.GuildMemberTimeout(
		message.GuildID,
		message.Author.ID,
		&timeout,
	); err != nil {
		h.Errors <- HandlerChannel{
			Message: "ProcessSpam(event): Unable to apply timeout to member",
			Fields: log.Fields{
				"user_id":       message.Author.ID,
				"heuristic_id":  rule.ID,
				"timeout":       timeout.String(),
				"error_message": err.Error(),
			},
		}

		return
	}

	bucket := make(map[string][]string)

	for idx := range logs {
		log := logs[idx]
		bucket[log.Message.ChannelID] = append(
			bucket[log.Message.ChannelID],
			log.Message.ID,
		)
	}

	h.Logs.ForEach(func(l *antispam.Log) {
		for idx := range logs {
			log := logs[idx]
			if l.Message.ID == log.Message.ID {
				l.MarkDeleted()
				break
			}
		}
	})

	var deleted int

	for channel, ids := range bucket {
		if err := session.ChannelMessagesBulkDelete(channel, ids); err == nil {
			deleted += len(ids)
		}
	}

	channels := make([]string, 0, len(bucket))
	for channel := range bucket {
		channels = append(channels, fmt.Sprintf("<#%s>", channel))
	}

	var fields []*discordgo.MessageEmbedField

	if rule.Duplicated {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name: "Content",
			Value: buildContentField(
				message.Content,
				message.Attachments,
				message.StickerItems,
			),
			Inline: true,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Channel(s)",
		Value:  strings.Join(channels, " "),
		Inline: true,
	})

	logUser(
		message.Author,
		log.WarnLevel,
		"ProcessSpam(event): Member timed out and messages deleted for triggering antispam",
		log.Fields{
			"deleted_count": deleted,
			"channel_count": len(channels),
			"heuristic_id":  rule.ID,
			"timeout":       rule.Timeout.String(),
		},
	)

	if deleted > 1 &&
		canViewChannel(session, message.GuildID, message.ChannelID) {
		message, err := sendSilentEmbed(session, h.Settings.LogChannel,
			&discordgo.MessageEmbed{
				Title: fmt.Sprintf(
					":shield: Spam detection triggered (%s)",
					strings.ToLower(rule.ID),
				),
				Description: fmt.Sprintf(
					"-# Attention: %d message(s) automatically removed from %d channel(s) due to suspected spam/phishing with potential malicious content. The user (%s) has been timed out for %s.",
					deleted,
					len(channels),
					message.Author.Mention(),
					rule.Timeout.String(),
				),
				Color: embedDeleteColor,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    message.Author.Username,
					IconURL: message.Author.AvatarURL("256"),
				},
				Fields: fields,
			},
		)
		if err != nil {
			h.Errors <- HandlerChannel{
				Message: "ProcessSpam(event): Unable to send message embed",
				Fields: log.Fields{
					"message_id":    message.ID,
					"error_message": err.Error(),
				},
			}

			return
		}

		h.ForwardAlert(session, message, true)
	}
}

func (h *Handler) ForwardAlert(
	session *discordgo.Session,
	message *discordgo.Message,
	ping bool,
) {
	if h.Settings.AlertChannel != "" {
		session.ChannelMessageSendReply(
			h.Settings.AlertChannel,
			"",
			message.Forward(),
		)

		if ping && h.Settings.ModRole != "" {
			session.ChannelMessageSend(
				h.Settings.AlertChannel,
				fmt.Sprintf("<@&%s>", h.Settings.ModRole),
			)
		}
	}
}

func (h *Handler) SendError(session *discordgo.Session, event HandlerChannel) {
	if h.Settings.AlertChannel != "" {
		fields := make([]*discordgo.MessageEmbedField, 0, len(event.Fields))
		for k, v := range event.Fields {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   k,
				Value:  TruncateContent(v.(string)),
				Inline: true,
			})
		}

		sendSilentEmbed(
			session,
			h.Settings.LogChannel,
			&discordgo.MessageEmbed{
				Title: ":no_entry: " + event.Message,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    session.State.User.Username,
					IconURL: session.State.User.AvatarURL("256"),
				},
				Fields: fields,
			},
		)

		log.WithFields(event.Fields).Error(event.Message)
	}
}
