// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/lcook/pulsar/internal/antispam"
	"github.com/lcook/pulsar/internal/cache"
	"github.com/lcook/pulsar/internal/config"
)

const (
	embedDeleteColor int = 0xDC322F
	embedUpdateColor int = 0x268BD2
)

type Handler struct {
	Settings config.Settings
	Events   []any
	Logs     *cache.RingBuffer[antispam.Log]
}

func New(settings config.Settings, buffer uint64) *Handler {
	h := &Handler{
		Settings: settings,
		Logs:     cache.NewRingBuffer[antispam.Log](buffer),
	}

	h.Events = append(h.Events, h.MessageCreate)
	h.Events = append(h.Events, h.MessageDelete)
	h.Events = append(h.Events, h.MessageUpdate)
	h.Events = append(h.Events, h.GuildMemberAdd)
	h.Events = append(h.Events, h.GuildMemberRemove)
	h.Events = append(h.Events, h.AutoModExecution)

	return h
}

func (h *Handler) ProcessSpam(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
	logs []*antispam.Log,
	rule *antispam.HeuristicRule,
) {
	timeout := time.Now().Add(rule.Timeout)
	session.GuildMemberTimeout(message.GuildID, message.Author.ID, &timeout)

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
			Name:   "Contents",
			Value:  buildContentField(message.Content, message.Attachments),
			Inline: true,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Channel(s)",
		Value:  strings.Join(channels, " "),
		Inline: true,
	})

	if deleted > 1 &&
		canViewChannel(session, message.GuildID, message.ChannelID) {
		session.ChannelMessageSendEmbed(
			h.Settings.LogChannel,
			&discordgo.MessageEmbed{
				Title: ":shield: Anti-spam alert",
				Description: fmt.Sprintf(
					"%d message(s) automatically removed from %d channel(s) due to suspected spam or advertising activity by <@%s>. The user has been timed out for %s. _Please exercise caution: these messages may contain malicious links, phishing attempts, or other harmful content_.",
					deleted,
					len(channels),
					message.Author.ID,
					rule.Timeout.String(),
				),
				Timestamp: time.Now().Format(time.RFC3339),
				Color:     embedDeleteColor,
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf(
						"ID: %s | HEURISTIC: %s",
						message.Author.ID,
						rule.ID,
					),
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    message.Author.Username,
					IconURL: message.Author.AvatarURL("256"),
				},
				Fields: fields,
			},
		)
	}
}
