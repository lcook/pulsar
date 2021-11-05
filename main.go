/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package main

import (
	"flag"
	"fmt"
	"net/http"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/bwmarrin/discordgo"
	"github.com/lcook/hookrelay"
	"github.com/lcook/pulseline/internal/config"
	"github.com/lcook/pulseline/internal/hook/git"
	log "github.com/sirupsen/logrus"
)

var (
	/*
	 * Default version displayed in the (-v) pulseline version
	 * command-line flag.  This is set during the build phase,
	 * possibly including the current git commit short-hash.
	 */
	Version = "devel"
)

type (
	Handler []hookrelay.Hook
)

func main() {
	var (
		verboseLevel int
		cfgFile      string
		version      bool
		color        bool
	)
	flag.IntVar(&verboseLevel, "V", 1, "Log verbosity level (1-3)")
	flag.StringVar(&cfgFile, "c", "config.yaml", "YAML configuration file path")
	flag.BoolVar(&version, "v", false, "Display pulseline version")
	flag.BoolVar(&color, "d", false, "Disable color output in logs")
	flag.Parse()
	log.SetFormatter(&nested.Formatter{
		ShowFullLevel:    true,
		NoUppercaseLevel: true,
		TrimMessages:     true,
		TimestampFormat:  "[02/Jan/2006:15:04:05]",
		NoFieldsColors:   true,
		NoColors:         color,
	})
	if version {
		fmt.Println(Version)
		return
	}
	/*
	 * Clamp the verbosity with an lower bound of 1 and
	 * upper bound of 3 (1-3).
	 */
	verboseClamp := func(n, lower, upper int) int {
		if n < lower {
			return lower
		}

		if n > upper {
			return upper
		}

		return n
	}(verboseLevel, 1, 3)

	switch verboseClamp {
	case 1:
		log.SetLevel(log.InfoLevel)
	case 2:
		log.SetLevel(log.DebugLevel)
	case 3:
		log.SetLevel(log.TraceLevel)
	}
	cfg, err := config.Load(cfgFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not open configuration file")
	}
	log.Infof("loaded configuration settings (%s)", cfgFile)

	log.Printf("init discord ...")
	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not create discord session")
	}
	err = dg.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not open discord connection")
	}
	log.WithFields(log.Fields{
		"id":   dg.State.User.ID,
		"user": dg.State.User.Username,
	}).Info("discord session started")

	_ = dg.UpdateGameStatus(0, Version)

	log.Printf("init pulseline-%s ...", Version)
	srv, err := hookrelay.InitMux(dg, Handler{
		&git.Pulse{
			Option: (hookrelay.DefaultOptions),
		},
	}, cfgFile, cfg.Port)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not start pulseline server")
	}

	log.WithFields(log.Fields{
		"port": cfg.Port,
	}).Info("pulseline server started")

	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not listen on port")
	}
	//nolint
	dg.Close()
}
