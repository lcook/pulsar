/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package bot

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/lcook/pulsar/internal/bot/command"
	"github.com/lcook/pulsar/internal/bot/event"
	"github.com/lcook/pulsar/internal/config"
	"github.com/lcook/pulsar/internal/pulse/hook/git"
	"github.com/lcook/pulsar/internal/relay"
	"github.com/lcook/pulsar/internal/version"
	log "github.com/sirupsen/logrus"
)

type Handler []relay.Hook

type Pulsar struct {
	config.Settings
}

type logError struct {
	*log.Entry
	Message string
	Error   error
}

func (p *Pulsar) Session(path string) (*discordgo.Session, *logError) {
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

	session.AddHandler(event.MessageDelete)
	session.AddHandler(event.MessageUpdate)

	session.Identify.Intents = discordgo.IntentsAll
	session.State.MaxMessageCount = 500

	_ = session.UpdateGameStatus(0, version.Build)

	log.Printf("init pulsar-%s ...", version.Build)

	srv, err := relay.InitMux(session, Handler{
		&git.Pulse{Option: (relay.DefaultOptions)},
	}, path, p.AcceptHost, p.AcceptPort)
	if err != nil {
		return nil, entry(err, "could not start pulsar server")
	}

	log.WithFields(log.Fields{
		"host": p.AcceptHost,
		"port": p.AcceptPort,
	}).Info("pulsar server started")

	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		return nil, entry(err, "could not listen on port")
	}

	return session, entry(nil, "")
}

func Run(path string) {
	pulsar, err := config.FromFile[Pulsar](path)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not open configuration file")
	}

	log.Infof("loaded configuration settings (%s)", path)

	session, entry := pulsar.Session(path)
	if entry.Error != nil {
		entry.Fatal(entry.Message)
	}

	_ = session.Close()
}
