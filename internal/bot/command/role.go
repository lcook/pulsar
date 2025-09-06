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
	roleRegex  string = rolePrefix + `\s(?P<` + roleSubExp + `>[A-z0-9\s]*)`
)

//go:embed data/roles.json
var roleData []byte

//go:embed templates/roles.tpl
var tplRoleData embed.FS

const tplRolePath string = "templates/roles.tpl"

func RoleHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID || message.Content == "" ||
		!strings.HasPrefix(message.Content, rolePrefix) {
		return
	}

	var roles map[string]string

	_ = json.Unmarshal(roleData, &roles)

	role := messageMatchRegex(message, roleRegex, roleSubExp)
	roleID := roles[role]

	if message.Content == rolePrefix ||
		strings.HasPrefix(message.Content, rolePrefix) && roleID == "" {
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

	if hasRole(message.Member, roleID) {
		err := session.GuildMemberRoleRemove(message.GuildID, message.Author.ID, roleID)
		if err != nil {
			_, _ = session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
				Title: fmt.Sprintf(
					"Error occurred when trying to remove role '%s' from %s",
					role,
					message.Author.GlobalName,
				),
				Author: &discordgo.MessageEmbedAuthor{
					Name:    message.Author.String(),
					IconURL: message.Author.AvatarURL("png"),
				},
				Color:       embedColor,
				Description: err.Error(),
			})

			return
		}

		_, _ = session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    message.Author.String(),
				IconURL: message.Author.AvatarURL("png"),
			},
			Color: embedColor,
			Description: fmt.Sprintf(
				"<@%s> was removed from the `%s` role.",
				message.Author.ID,
				role,
			),
		})
	} else {
		guildRoles, _ := session.GuildRoles(message.GuildID)
		for _, guildRole := range guildRoles {
			if guildRole.ID == roleID {
				err := session.GuildMemberRoleAdd(message.GuildID, message.Author.ID, roleID)
				if err != nil {
					_, _ = session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
						Title: fmt.Sprintf("Error occurred when trying to assign role '%s' to %s", role, message.Author.GlobalName),
						Author: &discordgo.MessageEmbedAuthor{
							Name:    message.Author.String(),
							IconURL: message.Author.AvatarURL("png"),
						},
						Color:       embedColor,
						Description: err.Error(),
					})

					return
				}

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
