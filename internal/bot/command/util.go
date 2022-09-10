/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2022, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package command

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func hasRole(member *discordgo.Member, id string) bool {
	for _, cid := range member.Roles {
		if cid == id {
			return true
		}
	}

	return false
}

func messageMatchRegex(message *discordgo.MessageCreate, regex, subexp string) string {
	reg := regexp.MustCompile(regex)
	if len(reg.FindStringSubmatch(message.Content)) == 2 {
		if match := reg.FindStringSubmatch(message.Content)[reg.SubexpIndex(subexp)]; match != "" {
			return match
		}
	}

	return ""
}
