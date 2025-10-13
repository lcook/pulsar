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

func (p *Pulsar) NewSession(path string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + p.Token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents = discordgo.IntentsAll
	session.State.MaxMessageCount = 500

	err = session.Open()
	if err != nil {
		return nil, err
	}

	log.Printf("Initialising pulsar-%s", version.Build)

	log.WithFields(log.Fields{
		"id":   session.State.User.ID,
		"user": session.State.User.Username,
	}).Info("Discord session started")

	session.UpdateGameStatus(0, version.Build)

	session.AddHandler(command.Bug.Handler)
	session.AddHandler(command.Role.Handler)
	session.AddHandler(command.Help.Handler)

	session.AddHandler(event.MessageDelete)
	session.AddHandler(event.MessageUpdate)
	session.AddHandler(event.GuildMemberRemove)
	session.AddHandler(event.AutoModExecution)

	hooks := Handler{
		&git.Pulse{Option: (relay.DefaultOptions)},
	}

	srv, err := relay.InitMux(session, hooks, path,
		p.AcceptHost, p.AcceptPort)
	if err != nil {
		return nil, err
	}

	for _, hook := range hooks {
		log.WithFields(log.Fields{
			"endpoint": hook.Endpoint(),
		}).Info("Registered mux handler")
	}

	log.WithFields(log.Fields{
		"host": p.AcceptHost,
		"port": p.AcceptPort,
	}).Infof("Initialised relay server with %d hook(s)", len(hooks))

	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		return nil, err
	}

	return session, err
}

func Run(path string) {
	pulsar, err := config.FromFile[Pulsar](path)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not open configuration file")
	}

	log.Infof("Loaded configuration settings (%s)", path)

	session, err := pulsar.NewSession(path)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not create pulsar session")
	}

	session.Close()
}
