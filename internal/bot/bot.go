package bot

import (
	"net/http"

	"github.com/bsdlabs/pulsar/internal/bot/command"
	"github.com/bsdlabs/pulsar/internal/pulse/hook/git"
	"github.com/bsdlabs/pulsar/internal/util"
	"github.com/bsdlabs/pulsar/internal/version"
	"github.com/bwmarrin/discordgo"
	"github.com/lcook/hookrelay"
	log "github.com/sirupsen/logrus"
)

type (
	Handler []hookrelay.Hook
)

func Run(config string) {
	cfg, err := util.GetConfig[Config](config)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not open configuration file")
	}

	log.Infof("loaded configuration settings (%s)", config)
	log.Printf("init discord ...")

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not create discord session")
	}

	err = session.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not open discord connection")
	}

	log.WithFields(log.Fields{
		"id":   session.State.User.ID,
		"user": session.State.User.Username,
	}).Info("discord session started")

	session.AddHandler(command.BugHandler)
	session.Identify.Intents = discordgo.IntentsGuildMessages

	_ = session.UpdateGameStatus(0, version.Build)

	log.Printf("init pulsar-%s ...", version.Build)

	srv, err := hookrelay.InitMux(session, Handler{
		&git.Pulse{
			Option: (hookrelay.DefaultOptions),
		},
	}, config, cfg.Port)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not start pulsar server")
	}

	log.WithFields(log.Fields{
		"port": cfg.Port,
	}).Info("pulsar server started")

	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not listen on port")
	}
	//nolint
	session.Close()
}
