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

	"github.com/bwmarrin/discordgo"
)

const (
	bugzBase  = "https://bugs.freebsd.org"
	bugz      = bugzBase + "/bugzilla/"
	bugzRest  = bugz + "/rest"
	bugzBug   = bugzRest + "/bug"
	bugzBugID = bugzBug + "?id=%s"
)

type bug struct {
	Status    string
	Summary   string
	Product   string
	Component string
	Creation  string `json:"creation_time"`
	Creator   struct {
		Email    string
		ID       int
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
	bugRegex := regexp.MustCompile(`bug\s!(?P<id>\d{1,6})`)
	if match := bugRegex.FindStringSubmatch(m.Content); len(match) == 2 {
		if id := bugRegex.FindStringSubmatch(m.Content)[bugRegex.SubexpIndex("id")]; id != "" {
			//nolint
			s.ChannelTyping(m.ChannelID)
			resp, err := http.Get(fmt.Sprintf(bugzBugID, id))
			if err != nil {
				//nolint
				s.ChannelMessageSend(m.ChannelID, "Could to fetch data from Bugzilla.")
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
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Could not find Bugzilla problem report with ID **%s**.", id))
				return
			}
			bug := pr.Bugs[0]
			embed := &discordgo.MessageEmbed{
				Color:       0x680000,
				Description: fmt.Sprintf("[%s](%s)", bug.Summary, fmt.Sprintf("%s/%s", bugzBase, id)),
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Status",
						Value:  bug.Status,
						Inline: true,
					},
					{
						Name:   "Product",
						Value:  bug.Product,
						Inline: true,
					},
					{
						Name:   "Component",
						Value:  bug.Component,
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%s • %s", id, bug.Creator.RealName),
				},
				Timestamp: bug.Creation,
			}
			//nolint
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		}
	}
}
