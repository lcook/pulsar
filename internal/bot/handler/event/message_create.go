// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/lcook/pulsar/internal/antispam"
)

func (h *Handler) MessageCreate(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Author.ID == s.State.User.ID || m.Member == nil || m.Author.Bot {
		return
	}

	if !h.Settings.Enabled {
		return
	}

	for _, id := range h.Settings.ExcludedRoleIDs {
		if slices.Contains(m.Member.Roles, id) {
			return
		}
	}

	var content strings.Builder
	for _, attachment := range m.Attachments {
		content.WriteString(attachment.Filename)
	}

	if m.Content != "" {
		content.WriteString(m.Content)
	}

	if content.String() == "" {
		return
	}

	sha := sha512.New()
	sha.Write([]byte(content.String()))

	hash := hex.EncodeToString(sha.Sum(nil))

	h.Logs.Add(antispam.Log{Message: m.Message, Hash: hash})

	if len(h.Logs.Slice()) > 1 {
		logs, rule := antispam.Run(m, hash, h.Logs, h.Settings.Rules)
		if len(logs) == 0 || rule == nil {
			return
		}

		timeout := time.Now().Add(rule.Timeout)
		s.GuildMemberTimeout(m.GuildID, m.Author.ID, &timeout)

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
			s.ChannelMessageSendEmbed(
				h.Settings.LogChannel,
				&discordgo.MessageEmbed{
					Title: ":shield: Anti-spam alert",
					Description: fmt.Sprintf(
						"%d message(s) automatically removed from %d channel(s) due to suspected spam or advertising activity by <@%s>. The user has been timed out for %s. _Please exercise caution: these messages may contain malicious links, phishing attempts, or other harmful content_.",
						deleted,
						len(channels),
						m.Author.ID,
						rule.Timeout.String(),
					),
					Timestamp: time.Now().Format(time.RFC3339),
					Color:     embedDeleteColor,
					Footer: &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf(
							"ID: %s | HEURISTIC: %s",
							m.Author.ID,
							rule.ID,
						),
					},
					Author: &discordgo.MessageEmbedAuthor{
						Name:    m.Author.Username,
						IconURL: m.Author.AvatarURL("256"),
					},
					Fields: fields,
				},
			)
		}
	}
}
