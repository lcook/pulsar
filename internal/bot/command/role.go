/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package command

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type roleAction struct {
	display string
	verb    string
	handler func(string, string, string, ...discordgo.RequestOption) error
}

const (
	rolePrefix string = prefix + "role"
	roleSubExp string = "roleid"
	roleRegex  string = rolePrefix + `\s(?P<` + roleSubExp + `>[A-z0-9-\s]*)`
)

//go:embed data/roles.json
var roleData []byte

var Role = Command{
	Name:        "role",
	Description: "Assign or remove a role to yourself",
	Handler: func(session *discordgo.Session, message *discordgo.MessageCreate) {
		if message.Author.ID == session.State.User.ID || message.Content == "" ||
			!strings.HasPrefix(message.Content, rolePrefix) {
			return
		}

		var roles map[string][]string
		err := json.Unmarshal(roleData, &roles)
		if err != nil {
			return
		}

		role := messageMatchRegex(message, roleRegex, roleSubExp)
		var roleID string
		if info, ok := roles[role]; ok && len(info) > 0 {
			roleID = info[0]
		}

		if role == "" || roleID == "" {
			fields := make([]*discordgo.MessageEmbedField, 0, len(roles))
			for _, info := range roles {
				fields = append(fields, &discordgo.MessageEmbedField{
					Value: fmt.Sprintf("<@&%s>: %s", info[0], info[1]),
				})
			}

			session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
				Title:       "Self-assignable roles",
				Description: "Assign yourself to any of the below roles by using _`!role <name>`_",
				Color:       embedColorFreeBSD,
				Fields:      fields,
			})

			return
		}

		guildRoles, _ := session.GuildRoles(message.GuildID)

		idx := slices.IndexFunc(guildRoles, func(r *discordgo.Role) bool {
			return r.ID == roleID
		})

		if idx < 0 {
			return
		}

		guildRole := guildRoles[idx]
		action := roleAction{"removed from", "remove", session.GuildMemberRoleRemove}
		if !hasRole(message.Member, roleID) {
			action = roleAction{"assigned to", "assign", session.GuildMemberRoleAdd}
		}

		if err := action.handler(message.GuildID, message.Author.ID, roleID); err != nil {
			session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("Error occurred when trying to %s <@&%s> <@!%s>", action.verb, role, message.Author.ID),
				Description: err.Error(),
				Color:       guildRole.Color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    message.Author.String(),
					IconURL: message.Author.AvatarURL("256"),
				},
			})

			return
		}

		session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
			Description: fmt.Sprintf("<@!%s> was %s the <@&%s> role", message.Author.ID, action.display, roleID),
			Color:       guildRole.Color,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    message.Author.String(),
				IconURL: message.Author.AvatarURL("256"),
			},
		})
	},
}
