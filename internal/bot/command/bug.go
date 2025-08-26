/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	bugzBase   string = "https://bugs.freebsd.org"
	bugz       string = bugzBase + "/bugzilla/"
	bugzRest   string = bugz + "/rest"
	bugzBug    string = bugzRest + "/bug"
	bugzBugID  string = bugzBug + "?id=%s"
	bugzReport string = bugzBase + "/%s"

	bugzRegexPrefix string = "#"
	bugzSubExp      string = "bugid"
	bugzRegex       string = `bug\s` + bugzRegexPrefix + `(?P<` + bugzSubExp + `>\d{1,6})`

	embedColor int = 0x680000
)

type bug struct {
	ID        json.Number
	Status    string
	Summary   string
	Product   string
	Component string
	Creation  string `json:"creation_time"`
	Creator   struct {
		Email    string
		ID       json.Number
		Name     string
		RealName string `json:"real_name"`
	} `json:"creator_detail"`
}

type problemReport struct {
	Bugs []bug `json:"bugs"`
}

func BugHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	if bugID := messageMatchRegex(message, bugzRegex, bugzSubExp); bugID != "" {
		_ = session.ChannelTyping(message.ChannelID)

		resp, err := http.Get(fmt.Sprintf(bugzBugID, bugID))
		if err != nil {
			_, _ = session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
				Title:       "FreeBSD Bugzilla",
				Color:       embedColor,
				Description: "Could not fetch data from Bugzilla.",
				Timestamp:   time.Now().Format(time.RFC3339),
			})

			return
		}
		defer resp.Body.Close()

		var report problemReport

		err = json.NewDecoder(resp.Body).Decode(&report)
		if err != nil {
			return
		}

		if len(report.Bugs) < 1 {
			_, _ = session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
				Title:       "FreeBSD Bugzilla",
				Color:       embedColor,
				Description: fmt.Sprintf("Could not find Bugzilla problem report with ID **%s**.", bugID),
				Timestamp:   time.Now().Format(time.RFC3339),
			})

			return
		}

		bug := &report.Bugs[0]
		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("FreeBSD Bugzilla - Bug %s", bug.ID),
			Color:       embedColor,
			Description: embedReport(bug),
			Footer: &discordgo.MessageEmbedFooter{
				Text: bug.Creator.RealName,
			},
			Timestamp: bug.Creation,
		}

		_, _ = session.ChannelMessageSendEmbed(message.ChannelID, embed)
	}
}
