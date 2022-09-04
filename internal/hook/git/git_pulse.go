/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021-2022, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import (
	"time"

	//nolint
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	repoPorts int = 0xB58900
	repoSrc   int = 0xDC322F
	repoDoc   int = 0x268BD2
)

type Pulse struct {
	Config struct {
		Endpoint string
		Secret   string

		WebhookID    string
		WebhookToken string
	} `yaml:"git"`
	Option byte
}

func (p *Pulse) Endpoint() string { return p.Config.Endpoint }
func (p *Pulse) Options() byte    { return p.Option }

func (p *Pulse) validHmac(buf []byte, w http.ResponseWriter, r *http.Request) bool {
	h := strings.SplitN(r.Header.Get("X-Hub-Signature"), "=", 2)
	if len(h) < 1 || h[0] != "sha1" {
		log.WithFields(log.Fields{
			"client": r.Header.Get("X-FORWARDED-FOR"),
		}).Warn(("git: X-Hub-Signature does not exist in request header"))
		return false
	}
	hm := hmac.New(sha1.New, []byte(p.Config.Secret))
	hm.Write(buf)
	eh := hex.EncodeToString(hm.Sum(nil))
	/*
	 * Make sure we contain a valid `X-Hub-Signature` header, as provided
	 * in the GitHub commit-payload.  Compute the HMAC hex digest with a
	 * locally stored secret (as defined within the configuration file) to
	 * ensure correct authenticity.
	 */
	if h[1] != eh {
		log.WithFields(log.Fields{
			"client": r.Header.Get("X-FORWARDED-FOR"),
		}).Warn(("git: unauthorized request received"))
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	return true
}

func (p *Pulse) Response(resp any) func(w http.ResponseWriter, r *http.Request) {
	dg := resp.(*discordgo.Session)
	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("git: failed to read payload")
			return
		}
		if !p.validHmac(buf, w, r) {
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
		/*
		 * Enumerate through all of the commits in the GitHub payload data,
		 * passing them off to a Discord Webhook that emits an embedded
		 * message containing relevant information of a commit.
		 */
		for idx, commit := range payload.Commits {
			log.WithFields(log.Fields{
				"commit":  commit.shortHash(),
				"author":  commit.Committer.String(),
				"message": strings.Split(commit.Message, "\n")[0],
			}).Trace("git: parsed commit")
			queue := fmt.Sprintf("%d/%d", idx+1, len(payload.Commits))

			params := &discordgo.WebhookParams{
				Username:  fmt.Sprintf("%s <%s@>", commit.Committer.Name, commit.Committer.Username),
				AvatarURL: commit.Committer.Avatar(),
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       color,
						Description: commit.embedCommit(payload.Repository.String(), payload.Ref),
						Footer: &discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%s repository", payload.Repository.String()),
						},
						Timestamp: commit.Timestamp.Format(time.RFC3339),
					},
				},
			}

			_, err = dg.WebhookExecute(p.Config.WebhookID, p.Config.WebhookToken, false, params)
			if err != nil {
				log.WithFields(log.Fields{
					"webhook": p.Config.WebhookID,
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
		//nolint
		defer r.Body.Close()
	}
}

func (p *Pulse) LoadConfig(config string) error {
	file, err := os.Open(config)
	if err != nil {
		return err
	}
	//nolint
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var cfg Pulse
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}
	p.Config.Secret = cfg.Config.Secret
	p.Config.Endpoint = cfg.Config.Endpoint
	p.Config.WebhookID = cfg.Config.WebhookID
	p.Config.WebhookToken = cfg.Config.WebhookToken
	return nil
}
