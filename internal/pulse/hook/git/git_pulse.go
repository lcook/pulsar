// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package git

import (
	"crypto/hmac"
	"crypto/sha1" //nolint
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"

	"github.com/lcook/pulsar/internal/config"
)

const (
	repoPorts int = 0xB58900
	repoSrc   int = 0xDC322F
	repoDoc   int = 0x268BD2
)

type Pulse struct {
	config.Settings
	Option byte
}

func (p *Pulse) Endpoint() string { return p.GithubWebhookEndpoint }

func (p *Pulse) Options() byte { return p.Option }

func (p *Pulse) validHmac(
	buf []byte,
	writer http.ResponseWriter,
	req *http.Request,
) bool {
	header := strings.SplitN(req.Header.Get("X-Hub-Signature"), "=", 2)
	if len(header) < 1 || header[0] != "sha1" {
		log.WithFields(log.Fields{
			"client": req.Header.Get("X-FORWARDED-FOR"),
		}).Warn(("git: X-Hub-Signature does not exist in request header"))

		return false
	}

	_hmac := hmac.New(sha1.New, []byte(p.GithubWebhookSecret))
	_hmac.Write(buf)
	hmacSum := hex.EncodeToString(_hmac.Sum(nil))
	// Make sure we contain a valid `X-Hub-Signature` header, as provided
	// in the GitHub commit-payload. Compute the HMAC hex digest with a
	// locally stored secret (as defined within the configuration file) to
	// ensure correct authenticity.
	if header[1] != hmacSum {
		log.WithFields(log.Fields{
			"client": req.Header.Get("X-FORWARDED-FOR"),
		}).Warn(("git: unauthorized request received"))
		writer.WriteHeader(http.StatusUnauthorized)

		return false
	}

	return true
}

func (p *Pulse) Response(
	resp any,
) func(w http.ResponseWriter, r *http.Request) {
	session := resp.(*discordgo.Session)

	return func(writer http.ResponseWriter, req *http.Request) {
		buf, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error("git: failed to read payload")
			return
		}

		if !p.validHmac(buf, writer, req) {
			return
		}

		payload, err := commitEventPayload(buf)
		if err != nil {
			log.Error("git: failed to unmarshal payload")
		}

		log.WithFields(log.Fields{
			"branch":     payload.Ref,
			"commits":    len(payload.Commits),
			"repository": payload.Repository,
		}).Debug("git: received github payload")

		var color int

		switch payload.Repository.String() {
		case "src":
			color = repoSrc
		case "ports":
			color = repoPorts
		case "doc":
			color = repoDoc
		}
		// Enumerate through all of the commits in the GitHub payload data,
		// passing them off to a Discord Webhook that emits an embedded
		// message containing relevant information of a commit.
		for idx, commit := range payload.Commits {
			log.WithFields(log.Fields{
				"commit":  commit.shortHash(),
				"author":  commit.Committer.String(),
				"message": strings.Split(commit.Message, "\n")[0],
			}).Trace("git: parsed commit")

			queue := fmt.Sprintf("%d/%d", idx+1, len(payload.Commits))

			params := &discordgo.WebhookParams{
				Username: commit.Committer.Name,
				AvatarURL: Avatar(
					commit.Committer.Username,
					commit.Committer.Email,
				),
				Embeds: []*discordgo.MessageEmbed{
					{
						Color: color,
						Description: commit.embedCommit(
							payload.Repository.String(),
							payload.Ref,
						),
						Footer: &discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf(
								"%s repository",
								payload.Repository.String(),
							),
						},
						Author: func() *discordgo.MessageEmbedAuthor {
							if commit.Committer.Name != commit.Author.Name {
								return &discordgo.MessageEmbedAuthor{
									Name: commit.Author.Name,
									IconURL: Avatar(
										commit.Author.Username,
										commit.Author.Email,
									),
								}
							}

							return &discordgo.MessageEmbedAuthor{}
						}(),
						Timestamp: commit.Timestamp.Format(time.RFC3339),
					},
				},
			}

			_, err = session.WebhookExecute(
				p.GithubWebhookID,
				p.GithubWebhookToken,
				false,
				params,
			)
			if err != nil {
				log.WithFields(log.Fields{
					"webhook": p.GithubWebhookID,
					"commit":  commit.shortHash(),
					"author":  commit.Committer.String(),
					"queue":   queue,
				}).Error("git: unable to send message")

				continue
			}

			log.WithFields(log.Fields{
				"commit": commit.shortHash(),
				"queue":  queue,
			}).Trace("git: sent message to discord")
		}

		defer req.Body.Close()
	}
}

func (p *Pulse) LoadConfig(path string) error {
	contents, err := config.FromFile[config.Settings](path)
	if err != nil {
		return err
	}

	p.Settings = contents

	return nil
}
