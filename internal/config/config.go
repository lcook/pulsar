package config

import (
	"github.com/BurntSushi/toml"
)

type BotSettings struct {
	Token         string `toml:"DiscordBotToken"`
	WebhookID     string `toml:"DiscordWebhookID"`
	WebhookToken  string `toml:"DiscordWebHookToken"`
	WebhookSecret string `toml:"GitHubWebhookSecret"`
	// Adhoc endpoint(s)
	GitEndpoint string `toml:"GitEndpoint"`
}

type ListenerSettings struct {
	AcceptHost string `toml:"SocketAcceptHost"`
	AcceptPort string `toml:"SocketAcceptPort"`
}

type Settings struct {
	BotSettings      `toml:"bot"`
	ListenerSettings `toml:"listener"`
}

func FromFile[T any](path string) (T, error) {
	var settings T

	_, err := toml.DecodeFile(path, &settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}
