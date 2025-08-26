/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package command

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lcook/pulsar/internal/util"
)

const (
	prefix string = "!"

	rolePrefix string = prefix + "role"
	roleSubExp string = "roleid"
	roleRegex  string = rolePrefix + `\s(?P<` + roleSubExp + `>[A-z0-9]*)`
)

//go:embed data/roles.json
var roleData []byte

//go:embed templates/roles.tpl
var tplRoleData embed.FS

const tplRolePath string = "templates/roles.tpl"

func RoleHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	var roles map[string]string

	_ = json.Unmarshal(roleData, &roles)

	if message.Content == rolePrefix {
		_, _ = session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
			Title: "Self-assignable roles",
			Description: util.EmbedDescription(tplRolePath, tplRoleData, map[string]any{
				"roles": func() string {
					var fields []string
					for name := range roles {
						fields = append(fields, fmt.Sprintf("`%s`", name))
					}

					return strings.Join(fields, " ")
				}(),
				"prefix": rolePrefix,
			}),
			Color: embedColor,
		})

		return
	}

	if role := messageMatchRegex(message, roleRegex, roleSubExp); role != "" {
		var roles map[string]string

		_ = json.Unmarshal(roleData, &roles)

		roleID := roles[role]
		if roleID == "" {
			return
		}

		if hasRole(message.Member, roleID) {
			_ = session.GuildMemberRoleRemove(message.GuildID, message.Author.ID, roleID)
			_, _ = session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{
					Name:    message.Author.String(),
					IconURL: message.Author.AvatarURL("png"),
				},
				Color:       embedColor,
				Description: fmt.Sprintf("<@%s> was removed from the `%s` role.", message.Author.ID, role),
			})
		} else {
			guildRoles, _ := session.GuildRoles(message.GuildID)
			for _, guildRole := range guildRoles {
				if guildRole.ID == roleID {
					_ = session.GuildMemberRoleAdd(message.GuildID, message.Author.ID, roleID)
					_, _ = session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
						Author: &discordgo.MessageEmbedAuthor{
							Name:    message.Author.String(),
							IconURL: message.Author.AvatarURL("png"),
						},
						Color:       embedColor,
						Description: fmt.Sprintf("<@%s> was given the `%s` role.", message.Author.ID, role),
					})

					break
				}
			}
		}
	}
}
