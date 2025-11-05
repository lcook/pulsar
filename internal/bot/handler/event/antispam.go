/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

func (h *Handler) AntiSpam(s *discordgo.Session, m *discordgo.MessageCreate, hash string) {
	var logs []Log

	for _, log := range h.Logs.Slice() {
		if m.Author.ID == log.Message.Author.ID && !log.deleted.Load() {
			logs = append(logs, log)
		}
	}

	var rules Heuristics

	err := yaml.Unmarshal(heuristicsData, &rules)
	if err != nil {
		return
	}

	spamLogs, rule := GetHeuristics(hash, logs, rules.Rules)
	if rule == nil || len(spamLogs) == 0 {
		return
	}

	timeout := time.Now().Add(rule.Timeout)
	s.GuildMemberTimeout(m.GuildID, m.Author.ID, &timeout)

	bucket := make(map[string][]string)
	for _, log := range spamLogs {
		bucket[log.Message.ChannelID] = append(bucket[log.Message.ChannelID], log.Message.ID)
	}

	h.Logs.ForEach(func(l *Log) {
		for _, log := range spamLogs {
			if l.Message.ID == log.Message.ID {
				l.deleted.Store(true)
				break
			}
		}
	})

	var deleted int

	for channel, ids := range bucket {
		if err := s.ChannelMessagesBulkDelete(channel, ids); err == nil {
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
			Value:  buildContentField(m.Content, m.Attachments),
			Inline: true,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Channel(s)",
		Value:  strings.Join(channels, " "),
		Inline: true,
	})

	if deleted > 1 && canViewChannel(s, m.GuildID, m.ChannelID) {
		s.ChannelMessageSendEmbed(h.Settings.LogChannel, &discordgo.MessageEmbed{
			Title: ":shield: Anti-spam alert",
			Description: fmt.Sprintf("%d message(s) automatically removed from %d channel(s) due to suspected spam or advertising activity by <@%s>. The user has been timed out for %s. _Please exercise caution: these messages may contain malicious links, phishing attempts, or other harmful content_.",
				deleted, len(channels), m.Author.ID, rule.Timeout.String()),
			Timestamp: time.Now().Format(time.RFC3339),
			Color:     embedDeleteColor,
			Footer:    &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s | HEURISTIC: %s", m.Author.ID, rule.ID)},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    m.Author.Username,
				IconURL: m.Author.AvatarURL("256"),
			},
			Fields: fields,
		})
	}
}
