/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package config

import (
	"time"

	"github.com/BurntSushi/toml"
)

type BotSettings struct {
	Token      string `toml:"DiscordBotToken"`
	Prefix     string `toml:"DiscordBotPrefix"`
	LogChannel string `toml:"DiscordLogChannelID"`
	//
	WebhookID     string `toml:"DiscordWebhookID"`
	WebhookToken  string `toml:"DiscordWebHookToken"`
	WebhookSecret string `toml:"GitHubWebhookSecret"`
	// Adhoc endpoint(s)
	GitEndpoint string `toml:"GitEndpoint"`
	//
	AntiSpamSettings `toml:"AntiSpam"`
}

type AntiSpamSettings struct {
	Enabled           bool          `toml:"Enabled"`
	MessageCacheSize  uint64        `toml:"MessageCacheSize"`
	ExcludedRoleIDs   []string      `toml:"ExcludeRoleIDs"`
	MinumumAccountAge time.Duration `toml:"MinimumAccountAge"`
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
