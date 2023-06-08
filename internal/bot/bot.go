/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021-2022, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package bot

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/lcook/hookrelay"
	"github.com/lcook/pulsar/internal/bot/command"
	"github.com/lcook/pulsar/internal/pulse/hook/git"
	"github.com/lcook/pulsar/internal/util"
	"github.com/lcook/pulsar/internal/version"
	log "github.com/sirupsen/logrus"
)

type Handler []hookrelay.Hook

type Pulsar struct {
	Port  string `json:"event_listener_port"`
	Token string `json:"discord_bot_token"`
}

type logError struct {
	*log.Entry
	Message string
	Error   error
}

func (p *Pulsar) Session(config string) (*discordgo.Session, *logError) {
	entry := func(err error, s string) *logError {
		return &logError{
			log.WithFields(log.Fields{
				"error": err,
			}), s, err,
		}
	}

	log.Printf("init discord session ...")

	session, err := discordgo.New("Bot " + p.Token)
	if err != nil {
		return nil, entry(err, "could not create discord session")
	}

	err = session.Open()
	if err != nil {
		return nil, entry(err, "could not open discord connection")
	}

	log.WithFields(log.Fields{
		"id":   session.State.User.ID,
		"user": session.State.User.Username,
	}).Info("discord session started")

	session.AddHandler(command.BugHandler)
	session.AddHandler(command.RoleHandler)

	session.Identify.Intents = discordgo.IntentsGuildMessages

	_ = session.UpdateGameStatus(0, version.Build)

	log.Printf("init pulsar-%s ...", version.Build)

	srv, err := hookrelay.InitMux(session, Handler{
		&git.Pulse{Option: (hookrelay.DefaultOptions)},
	}, config, p.Port)

	if err != nil {
		return nil, entry(err, "could not start pulsar server")
	}

	log.WithFields(log.Fields{
		"port": p.Port,
	}).Info("pulsar server started")

	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		return nil, entry(err, "could not listen on port")
	}

	return session, entry(nil, "")
}

func Run(config string) {
	pulsar, err := util.GetConfig[Pulsar](config)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not open configuration file")
	}

	log.Infof("loaded configuration settings (%s)", config)

	session, entry := pulsar.Session(config)
	if entry.Error != nil {
		entry.Fatal(entry.Message)
	}

	_ = session.Close()
}
