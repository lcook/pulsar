// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package command

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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
		Title:       ":robot: pulsar bot status",
		Description: "-# Source code available on [GitHub](https://github.com/lcook/pulsar).",
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
				Name:   "Latency",
				Value:  s.HeartbeatLatency().String(),
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
				Name: fmt.Sprintf("Commands (%d)", len(h.commands)),
				Value: func() string {
					var buf strings.Builder
					for _, command := range h.commands {
						fmt.Fprintf(
							&buf,
							" `%s%s`",
							h.Settings.Prefix,
							command.Name,
						)
					}

					return buf.String()
				}(),
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
				Inline: true,
			},
		},
	})
}
