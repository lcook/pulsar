/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package bot

import (
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/lcook/pulsar/internal/config"
	"github.com/lcook/pulsar/internal/version"
)

type Bot struct {
	Settings config.Settings
	Session  *discordgo.Session
}

func New(path string) (*Bot, error) {
	log.WithFields(log.Fields{
		"file": path,
	}).Info("Loading configuration settings")

	settings, err := config.FromFile[config.Settings](path)
	if err != nil {
		return nil, err
	}

	log.Info("Setting up new Discord session")

	session, err := discordgo.New("Bot " + settings.Token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Settings: settings,
		Session:  session,
	}, nil
}

func (b *Bot) Run(handlers ...[]any) error {
	b.Session.Identify.Intents = discordgo.IntentsAll
	b.Session.State.MaxMessageCount = 500

	log.Info("Starting websocket connection with Discord")

	err := b.Session.Open()
	if err != nil {
		return err
	}

	b.Session.UpdateGameStatus(0, version.Build)

	var _handlers []any

	for _, slice := range handlers {
		_handlers = append(_handlers, slice...)
	}

	log.Infof("Registering %d Discord event handler(s)", len(_handlers))

	for _, handler := range _handlers {
		b.Session.AddHandler(handler)
	}

	return nil
}
