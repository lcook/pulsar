/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package event

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

func (h *Handler) AntiSpam(s *discordgo.Session, m *discordgo.MessageCreate, hash string) {
	now := time.Now().UTC()

	var logs []Log
	for _, log := range h.Logs.Slice() {
		if m.Author.ID == log.Message.Author.ID {
			if now.Sub(log.Message.Timestamp.UTC()) > h.Settings.MessageWindow || log.deleted.Load() {
				continue
			}
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

	var deleted int

	h.Logs.ForEach(func(l *Log) {
		if l.deleted.Load() {
			return
		}

		for _, log := range spamLogs {
			if l.Message.ID == log.Message.ID {
				l.deleted.Store(true)

				if err := s.ChannelMessageDelete(l.Message.ChannelID, l.Message.ID); err == nil {
					deleted++
				}

				return
			}
		}
	})

	channels := make([]string, 0, len(spamLogs))
	for _, log := range spamLogs {
		channels = append(channels, fmt.Sprintf("<#%s>", log.Message.ChannelID))
	}

	channels = slices.Compact(channels)

	timeout := time.Now().Add(rule.Timeout)
	s.GuildMemberTimeout(m.GuildID, m.Author.ID, &timeout)

	var fields []*discordgo.MessageEmbedField

	if rule.Duplicated {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Contents",
			Value: buildContentField(m.Content, m.Attachments),
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "Channel(s)",
		Value: strings.Join(channels, " "),
	})

	if deleted > 1 && canViewChannel(s, m.GuildID, m.ChannelID) {
		s.ChannelMessageSendEmbed(h.Settings.LogChannel, &discordgo.MessageEmbed{
			Title: ":shield: Anti-spam alert",
			Description: fmt.Sprintf("%d message(s) automatically removed from %d channel(s) due to suspected spam or advertising activity by <@%s>. The user has been timed out for %s. _Please exercise caution: these messages may contain malicious links, phishing attempts, or other harmful content_",
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
