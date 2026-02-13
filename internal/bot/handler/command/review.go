// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/lcook/pulsar/internal/bot/handler/event"
)

const (
	reviewBase        string = "https://reviews.freebsd.org"
	conduitAPI        string = reviewBase + "/api/"
	conduitEmbedColor int    = 0x4a5f88
)

type Differential struct {
	Result struct {
		Data []struct {
			ID     int `json:"id"`
			Fields struct {
				Title  string `json:"title"`
				Author string `json:"authorPHID"`
				Status struct {
					Name string `json:"name"`
				} `json:"status"`
				Summary string `json:"summary"`
				Created int    `json:"dateCreated"`
			} `json:"fields"`
		} `json:"data"`
	} `json:"result"`
}

type DifferentialUser struct {
	Result struct {
		Data []struct {
			ID     int    `json:"id"`
			PHID   string `json:"phid"`
			Fields struct {
				Username string `json:"username"`
				Realname string `json:"realName"`
			} `json:"fields"`
		}
	}
}

type ConduitClient struct {
	Token string
}

func NewConduit(token string) *ConduitClient {
	return &ConduitClient{token}
}

func (c *ConduitClient) Call(
	endpoint string,
	fields map[string]string,
) (string, error) {
	values := url.Values{}
	values.Set("api.token", c.Token)

	for key, value := range fields {
		values.Set(key, value)
	}

	resp, err := http.PostForm(conduitAPI+endpoint, values)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *ConduitClient) GetUser(phid string) (DifferentialUser, error) {
	data, err := c.Call(
		"user.search",
		map[string]string{"constraints[phids][]": phid},
	)
	if err != nil {
		return DifferentialUser{}, err
	}

	var tmp DifferentialUser

	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		return DifferentialUser{}, err
	}

	return tmp, nil
}

func (c *ConduitClient) GetDifferential(id string) (Differential, error) {
	data, err := c.Call(
		"differential.revision.search",
		map[string]string{"constraints[ids][]": id},
	)
	if err != nil {
		return Differential{}, err
	}

	var tmp Differential

	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		return Differential{}, err
	}

	return tmp, nil
}

func (h *Handler) Review(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}

	diffRegex := `(?:(?i)\` + h.Settings.Prefix + `review\s+(?:D)?|(?:https?://)?reviews\.freebsd\.org/D)(?P<diffid>\d+)`
	if diffID := messageMatchRegex(m, diffRegex, "diffid"); diffID != "" {
		s.ChannelTyping(m.ChannelID)

		conduit := NewConduit(h.Settings.ConduitToken)

		author := &discordgo.MessageEmbedAuthor{
			Name:    "Phabricator: Differential D" + diffID,
			IconURL: "https://reviews.freebsd.org/file/data/qlge5ptgqas6r46gkigm/PHID-FILE-xrbh6ayr3mccyyz5tu5n/favicon",
		}

		data, err := conduit.GetDifferential(diffID)
		if err != nil {
			s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
				Description: fmt.Sprintf(
					"Unable to request data from Phabricator: %v",
					err,
				),
				Color:  conduitEmbedColor,
				Author: author,
			}, m.Reference())

			return
		}

		if len(data.Result.Data) < 1 {
			s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
				Description: fmt.Sprintf(
					"Unable to find Differential revision with ID matching **D%s**",
					diffID,
				),
				Color:  conduitEmbedColor,
				Author: author,
			}, m.Reference())

			return
		}

		diff := data.Result.Data[0]

		var fields []*discordgo.MessageEmbedField

		if diff.Fields.Status.Name != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Status",
				Value:  diff.Fields.Status.Name,
				Inline: true,
			})
		}

		if user, err := conduit.GetUser(diff.Fields.Author); err == nil {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name: "Author",
				Value: fmt.Sprintf(
					"%s <%s>",
					user.Result.Data[0].Fields.Realname,
					user.Result.Data[0].Fields.Username,
				),
				Inline: true,
			})
		}

		if diff.Fields.Summary != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "Summary",
				Value: event.TruncateContent(diff.Fields.Summary),
			})
		}

		created := time.Unix(int64(diff.Fields.Created), 0)

		s.ChannelMessageSendEmbedReply(
			m.ChannelID,
			&discordgo.MessageEmbed{
				Description: fmt.Sprintf(
					"[%s](%s)",
					diff.Fields.Title,
					fmt.Sprintf("%s/D%s", reviewBase, diffID),
				),
				Timestamp: created.Format(time.RFC3339),
				Color:     conduitEmbedColor,
				Author:    author,
				Fields:    fields,
			},
			m.Reference(),
		)
	}
}
