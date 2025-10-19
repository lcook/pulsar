/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package main

import (
	"flag"
	"fmt"
	"net/http"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"

	"github.com/lcook/pulsar/internal/bot"
	"github.com/lcook/pulsar/internal/bot/handler/command"
	"github.com/lcook/pulsar/internal/bot/handler/event"
	"github.com/lcook/pulsar/internal/pulse/hook/git"
	"github.com/lcook/pulsar/internal/relay"
	"github.com/lcook/pulsar/internal/version"
)

func main() {
	var (
		cfgFile     string
		color       bool
		showVersion bool
		verbosity   int
	)

	flag.IntVar(&verbosity, "V", 1, "Log verbosity level (1-3)")
	flag.StringVar(&cfgFile, "c", "config.toml", "TOML configuration file path")
	flag.BoolVar(&showVersion, "v", false, "Display pulsar version")
	flag.BoolVar(&color, "d", false, "Disable color output in logs")
	flag.Parse()

	log.SetFormatter(&nested.Formatter{
		ShowFullLevel:   true,
		TrimMessages:    true,
		TimestampFormat: "[02/Jan/2006:15:04:05]",
		NoFieldsColors:  true,
		NoColors:        color,
	})

	if showVersion {
		fmt.Println(version.Build)
		return
	}
	/*
	 * Clamp the verbosity with an lower bound of 1 and
	 * upper bound of 3 (1-3).
	 */
	if verbosity < 1 {
		verbosity = 1
	}

	if verbosity > 3 {
		verbosity = 3
	}

	switch verbosity {
	case 1:
		log.SetLevel(log.InfoLevel)
	case 2:
		log.SetLevel(log.DebugLevel)
	case 3:
		log.SetLevel(log.TraceLevel)
	}

	pulsar, err := bot.New(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Initialising pulsar-%s", version.Build)

	err = pulsar.Run(
		command.New(pulsar.Settings).Handlers(),
		event.New(pulsar.Settings).Events,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"id":   pulsar.Session.State.User.ID,
		"user": pulsar.Session.State.User.Username,
	}).Info("Discord session started")

	hooks := []relay.Hook{
		&git.Pulse{Option: (relay.DefaultOptions)},
	}

	srv, err := relay.InitMux(pulsar.Session, hooks, cfgFile,
		pulsar.Settings.AcceptHost, pulsar.Settings.AcceptPort)
	if err != nil {
		log.Fatal(err)
	}

	for _, hook := range hooks {
		log.WithFields(log.Fields{
			"endpoint": hook.Endpoint(),
		}).Info("Registered mux handler")
	}

	log.WithFields(log.Fields{
		"host": pulsar.Settings.AcceptHost,
		"port": pulsar.Settings.AcceptPort,
	}).Infof("Initialised relay server with %d hook(s)", len(hooks))

	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatal(err)
	}

	pulsar.Session.Close()
}
