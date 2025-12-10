// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"

	"github.com/lcook/pulsar/internal/bot"
	"github.com/lcook/pulsar/internal/bot/handler/command"
	"github.com/lcook/pulsar/internal/bot/handler/event"
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
	flag.StringVar(&cfgFile, "c", "config.yaml", "YAML configuration file path")
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

reload:
	pulsar, err := bot.New(cfgFile)

	if err != nil {
		log.Fatal(err)
	}

	identifier := fmt.Sprintf("pulsar-bot-%s", version.Build)

	err = pulsar.Init(
		identifier,
		discordgo.IntentGuilds|
			discordgo.IntentGuildMembers|
			discordgo.IntentGuildModeration|
			discordgo.IntentGuildMessages|
			discordgo.IntentAutoModerationExecution,
		true,
		command.New(pulsar.Settings).Handlers(),
		event.New(pulsar.Settings, pulsar.Settings.MessageCacheSize).Events,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Info(
		identifier + fmt.Sprintf(
			" is now running with PID=%d. Press CTRL-C to exit",
			os.Getpid(),
		),
	)

	sc := make(chan os.Signal, 1)
	signal.Notify(
		sc,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGUSR2,
	)

	switch <-sc {
	case syscall.SIGUSR2:
		log.Warn("SIGUSR signal received, reloading")
		goto reload
	case os.Interrupt, syscall.SIGINT, syscall.SIGTERM:
		log.Warn("Terminating signal received, closing down session")
	}

	err = pulsar.Session.Close()
	if err != nil {
		log.Error("could not close session gracefully")
	}
}
