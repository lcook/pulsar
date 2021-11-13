/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
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
	bugzRegex  string = `bug\s!(?P<id>\d{1,6})`

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

func BugHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	reg := regexp.MustCompile(bugzRegex)
	if match := reg.FindStringSubmatch(m.Content); len(match) == 2 {
		if id := reg.FindStringSubmatch(m.Content)[reg.SubexpIndex("id")]; id != "" {
			//nolint
			s.ChannelTyping(m.ChannelID)
			resp, err := http.Get(fmt.Sprintf(bugzBugID, id))
			if err != nil {
				//nolint
				s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title:       "FreeBSD Bugzilla",
					Color:       embedColor,
					Description: "Could not fetch data from Bugzilla.",
					Timestamp:   time.Now().Format(time.RFC3339),
				})
				return
			}
			//nolint
			defer resp.Body.Close()
			var pr problemReport
			err = json.NewDecoder(resp.Body).Decode(&pr)
			if err != nil {
				return
			}
			if len(pr.Bugs) < 1 {
				//nolint
				s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title:       "FreeBSD Bugzilla",
					Color:       embedColor,
					Description: fmt.Sprintf("Could not find Bugzilla problem report with ID **%s**.", id),
					Timestamp:   time.Now().Format(time.RFC3339),
				})
				return
			}
			bug := pr.Bugs[0]
			embed := &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("FreeBSD Bugzilla - Bug %s", bug.ID),
				Color:       embedColor,
				Description: embedDescription(bug),
				Footer: &discordgo.MessageEmbedFooter{
					Text: bug.Creator.RealName,
				},
				Timestamp: bug.Creation,
			}
			//nolint
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		}
	}
}
