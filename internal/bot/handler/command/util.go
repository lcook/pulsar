// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package command

import (
	"regexp"
	"slices"

	"github.com/bwmarrin/discordgo"
)

func hasRole(member *discordgo.Member, id string) bool {
	return slices.Contains(member.Roles, id)
}

func messageMatchRegex(
	message *discordgo.MessageCreate,
	regex, subexp string,
) string {
	reg := regexp.MustCompile(regex)
	if len(reg.FindStringSubmatch(message.Content)) == 2 {
		if match := reg.FindStringSubmatch(message.Content)[reg.SubexpIndex(subexp)]; match != "" {
			return match
		}
	}

	return ""
}

//nolint:unused
func directMessage(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
) bool {
	channel, _ := session.Channel(message.ChannelID)

	return channel.Type == discordgo.ChannelTypeDM
}
