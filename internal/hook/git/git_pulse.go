/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import (
	"crypto/hmac"
	//nolint
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	middlewareHmac          string = "hmac"
	middlewareDiscordEmebed string = "embed"
)

type Pulse struct {
	Config struct {
		Secret     string
		Endpoint   string
		Channel    string
		RepoSrc    string `yaml:"src"`
		RepoPorts  string `yaml:"ports"`
		RepoDoc    string `yaml:"doc"`
		Middleware []string
	} `yaml:"git"`
	Option byte
}

func (p *Pulse) Endpoint() string { return p.Config.Endpoint }
func (p *Pulse) Options() byte    { return p.Option }

func (p *Pulse) hasMiddleware(m string) bool {
	for _, v := range p.Config.Middleware {
		if v == m {
			return true
		}
	}
	return false
}

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

func (p *Pulse) Response(r interface{}) func(w http.ResponseWriter, r *http.Request) {
	dg := r.(*discordgo.Session)
	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("git: failed to read payload")
			return
		}
		if p.hasMiddleware(middlewareHmac) &&
			!p.validHmac(buf, w, r) {
			return
		}
		pl, err := commitEventPayload(buf)
		if err != nil {
			log.Error("git: failed to unmarshal payload")
		}
		log.WithFields(log.Fields{
			"branch":     pl.Ref,
			"commits":    len(pl.Commits),
			"repository": pl.Repository.Name,
		}).Debug("git: received github payload")
		if p.hasMiddleware(middlewareDiscordEmebed) {
			var color int
			switch pl.Repository.Name {
			case "src":
				color, _ = strconv.Atoi(p.Config.RepoSrc)
			case "ports":
				color, _ = strconv.Atoi(p.Config.RepoPorts)
			case "doc":
				color, _ = strconv.Atoi(p.Config.RepoDoc)
			}
			/*
			 * Iterate through each of the commits in the payload data, which
			 * are then sent as a Discord embedded message to a defined channel.
			 */
			for i, c := range pl.Commits {
				log.WithFields(log.Fields{
					"commit":  c.shortHash(),
					"author":  c.Committer.Name,
					"message": strings.Split(c.Message, "\n")[0],
				}).Trace("git: parsed commit")
				queue := fmt.Sprintf("%d/%d", i+1, len(pl.Commits))
				_, err = dg.ChannelMessageSendEmbed(p.Config.Channel, c.embed(pl.Repository.Name, pl.Ref, color))
				if err != nil {
					log.WithFields(log.Fields{
						"channel": p.Config.Channel,
						"commit":  c.shortHash(),
						"author":  c.Committer.Name,
						"queue":   queue,
					}).Error("git: unable to send message")
					continue
				}
				log.WithFields(log.Fields{
					"channel": p.Config.Channel,
					"commit":  c.shortHash(),
					"queue":   queue,
				}).Trace("git: sent message to discord")
			}
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
	hexToInt := func(hex string) int {
		res, err := strconv.ParseInt(strings.TrimPrefix(hex, "#"), 16, 64)
		if err != nil {
			return 000000
		}
		return int(res)
	}
	p.Config.Secret = cfg.Config.Secret
	p.Config.Endpoint = cfg.Config.Endpoint
	p.Config.Channel = cfg.Config.Channel
	p.Config.RepoSrc = strconv.Itoa(hexToInt(cfg.Config.RepoSrc))
	p.Config.RepoPorts = strconv.Itoa(hexToInt(cfg.Config.RepoPorts))
	p.Config.RepoDoc = strconv.Itoa(hexToInt(cfg.Config.RepoDoc))
	p.Config.Middleware = cfg.Config.Middleware
	return nil
}
