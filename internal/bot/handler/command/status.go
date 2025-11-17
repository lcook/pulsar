// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package command

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/lcook/pulsar/internal/version"
)

func (h *Handler) Status(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != h.Settings.Prefix+"status" {
		return
	}

	if directMessage(s, m) {
		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       ":ping_pong: pulsar bot",
		Description: "Source code available on [GitHub](https://github.com/lcook/pulsar).",
		Color:       embedColorFreeBSD,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Build",
				Value: version.Build,
			},
			{
				Name:   "Uptime",
				Value:  time.Since(h.Started).String(),
				Inline: true,
			},
			{
				Name: "Hostname",
				Value: func() string {
					hostname, _ := os.Hostname()
					return hostname
				}(),
				Inline: true,
			},
			{
				Name:   "Platform",
				Value:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
				Inline: true,
			},
			{
				Name: "Antispam",
				Value: func() string {
					if !h.Settings.Enabled {
						return "disabled"
					}

					return "enabled"
				}(),
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
