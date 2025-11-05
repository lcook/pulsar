/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package command

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	roleSubExp string = "roleid"
)

type roleAction struct {
	display string
	verb    string
	handler func(string, string, string, ...discordgo.RequestOption) error
}

func (h *Handler) Role(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}

	var (
		rolePrefix = h.Settings.Prefix + "role"
		roleRegex  = rolePrefix + `\s(?P<` + roleSubExp + `>[A-z0-9-\s]*)`
	)

	if m.Content == "" || !strings.HasPrefix(m.Content, rolePrefix) {
		return
	}

	roles := h.Settings.Roles
	role := messageMatchRegex(m, roleRegex, roleSubExp)

	var roleID string
	if info, ok := h.Settings.Roles[role]; ok {
		roleID = info.ID
	}

	if role == "" || roleID == "" {
		fields := make([]*discordgo.MessageEmbedField, 0, len(roles))
		for _, info := range roles {
			fields = append(fields, &discordgo.MessageEmbedField{
				Value: fmt.Sprintf("<@&%s>: %s", info.ID, info.Description),
			})
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "Self-assignable roles",
			Description: "Assign yourself to any of the below roles by using _`!role <name>`_.",
			Color:       embedColorFreeBSD,
			Fields:      fields,
		})

		return
	}

	guildRoles, _ := s.GuildRoles(m.GuildID)

	idx := slices.IndexFunc(guildRoles, func(r *discordgo.Role) bool {
		return r.ID == roleID
	})

	if idx < 0 {
		return
	}

	guildRole := guildRoles[idx]

	action := roleAction{"removed from", "remove", s.GuildMemberRoleRemove}
	if !hasRole(m.Member, roleID) {
		action = roleAction{"assigned to", "assign", s.GuildMemberRoleAdd}
	}

	if err := action.handler(m.GuildID, m.Author.ID, roleID); err != nil {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Error occurred when trying to %s <@&%s> <@!%s>", action.verb, role, m.Author.ID),
			Description: err.Error(),
			Color:       guildRole.Color,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    m.Author.String(),
				IconURL: m.Author.AvatarURL("256"),
			},
		})

		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("<@!%s> was %s the <@&%s> role", m.Author.ID, action.display, roleID),
		Color:       guildRole.Color,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.String(),
			IconURL: m.Author.AvatarURL("256"),
		},
	})
}
