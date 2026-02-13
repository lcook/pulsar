// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package command

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) User(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, h.Settings.Prefix+"user") {
		return
	}

	id := strings.Split(m.Content, " ")
	if len(id) == 1 {
		return
	}

	user, err := s.User(id[1])
	if err != nil {
		return
	}

	var fields []*discordgo.MessageEmbedField

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Username",
		Value:  user.Username,
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Global name",
		Value:  user.GlobalName,
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name: "Account created",
		Value: func() string {
			created, _ := discordgo.SnowflakeTimestamp(user.ID)
			return created.String()
		}(),
	})

	member, err := s.GuildMember(m.GuildID, user.ID)
	if err == nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Member joined",
			Value: member.JoinedAt.String(),
		})

		if member.Nick != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Member nickname",
				Value:  member.Nick,
				Inline: true,
			})
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name: "Pending",
			Value: func() string {
				if member.Pending {
					return "true"
				}

				return "false"
			}(),
			Inline: true,
		})

		guildRoles, err := s.GuildRoles(m.GuildID)
		if err != nil {
			return
		}

		rolesMap := make(map[string]*discordgo.Role, len(guildRoles))
		for _, role := range guildRoles {
			rolesMap[role.ID] = role
		}

		var roles []string

		for _, roleID := range member.Roles {
			if role, exists := rolesMap[roleID]; exists {
				roles = append(roles, role.Mention())
			}
		}

		if len(roles) >= 1 {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Roles",
				Value:  strings.Join(roles, " "),
				Inline: true,
			})
		}
	}

	s.ChannelMessageSendEmbedReply(
		m.ChannelID,
		&discordgo.MessageEmbed{
			Fields: fields,
			Image: &discordgo.MessageEmbedImage{
				URL: user.BannerURL("512"),
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: user.AvatarURL("512"),
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "ID: " + user.ID,
			},
			Color: user.AccentColor,
		},
		m.Reference(),
	)
}
