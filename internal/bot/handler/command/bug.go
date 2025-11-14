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
	bugzBase  string = "https://bugs.freebsd.org"
	bugz      string = bugzBase + "/bugzilla/"
	bugzRest  string = bugz + "/rest"
	bugzBug   string = bugzRest + "/bug"
	bugzBugID string = bugzBug + "?id=%s"

	bugzSubExp string = "bugid"
)

type user struct {
	Email    string
	ID       json.Number
	Nam      string
	RealName string `json:"real_name"`
}

type bug struct {
	ID         json.Number
	Status     string
	Resolution string
	Summary    string
	Product    string
	Component  string
	Version    string
	Platform   string
	Assignee   user   `json:"assigned_to_detail"`
	Creation   string `json:"creation_time"`
	Creator    user   `json:"creator_detail"`
}

type report struct {
	Bugs []bug `json:"bugs"`
}

func getReport(id string) (report, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(bugzBugID, id), nil)
	if err != nil {
		return report{}, err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return report{}, err
	}
	defer resp.Body.Close()

	var rep report

	err = json.NewDecoder(resp.Body).Decode(&rep)
	if err != nil {
		return report{}, err
	}

	return rep, nil
}

func (h *Handler) Bug(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}

	bugzRegex := `(?:(?i)\` + h.Settings.Prefix + `bug\s+|(?:https?://)?bugs\.freebsd\.org/bugzilla/show_bug\.cgi\?id=)(?P<bugid>\d{1,6})`
	if bugID := messageMatchRegex(m, bugzRegex, bugzSubExp); bugID != "" {
		s.ChannelTyping(m.ChannelID)

		author := &discordgo.MessageEmbedAuthor{
			Name:    "Bugzilla - report #" + bugID,
			IconURL: "https://raw.githubusercontent.com/freebsd/freebsd-src/refs/heads/main/stand/images/freebsd-logo-rev.png",
		}

		report, err := getReport(bugID)
		if err != nil {
			s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
				Description: fmt.Sprintf(
					"Unable to request data from Bugzilla: %v",
					err,
				),
				Timestamp: time.Now().Format(time.RFC3339),
				Color:     embedColorFreeBSD,
				Author:    author,
			}, m.Reference())

			return
		}

		if len(report.Bugs) < 1 {
			s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
				Description: fmt.Sprintf(
					"Unable to find Bugzilla report with ID matching **%s**",
					bugID,
				),
				Timestamp: time.Now().Format(time.RFC3339),
				Color:     embedColorFreeBSD,
				Author:    author,
			}, m.Reference())

			return
		}

		bug := &report.Bugs[0]

		s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
			Description: fmt.Sprintf(
				"[%s](%s/%s)",
				bug.Summary,
				bugzBase,
				bugID,
			),
			Timestamp: bug.Creation,
			Color:     embedColorFreeBSD,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf(
					"%s <%s>",
					bug.Creator.RealName,
					bug.Creator.Email,
				),
			},
			Author: author,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: "Status",
					Value: func() string {
						if bug.Resolution == "" {
							return bug.Status
						}

						return bug.Status + " " + bug.Resolution
					}(),
					Inline: true,
				},
				{Name: "Product", Value: bug.Product, Inline: true},
				{Name: "Component", Value: bug.Component, Inline: true},
				{Name: "Version", Value: bug.Version, Inline: true},
				{Name: "Platform", Value: bug.Platform, Inline: true},
				{
					Name: "Assignee",
					Value: func() string {
						if bug.Assignee.RealName == "" {
							return "Nobody"
						}

						return bug.Assignee.RealName
					}(),
					Inline: true,
				},
			},
		}, m.Reference())
	}
}
