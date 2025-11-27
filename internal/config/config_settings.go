// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package config

import (
	"time"

	"github.com/lcook/pulsar/internal/antispam"
)

type BotSettings struct {
	Token        string   `yaml:"discord_token"`
	Prefix       string   `yaml:"discord_prefix"`
	Commands     []string `yaml:"discord_commands"`
	LogChannel   string   `yaml:"discord_log_channel_id"`
	AlertChannel string   `yaml:"discord_alert_channel_id"`
	ModRole      string   `yaml:"discord_mod_role_id"`

	GithubWebhookID    string `yaml:"discord_github_webhook_id"`
	GithubWebhookToken string `yaml:"discord_github_webhook_token"`

	ConduitToken string `yaml:"discord_conduit_token"`

	Roles map[string]Role `yaml:"roles"`

	AntiSpamSettings `yaml:"antispam"`
}

type AntiSpamSettings struct {
	Enabled           bool                     `yaml:"enabled"`
	MessageCacheSize  uint64                   `yaml:"message_cache_size"`
	ExcludedRoleIDs   []string                 `yaml:"excluded_role_ids"`
	MinumumAccountAge time.Duration            `yaml:"minimum_account_age"`
	Rules             []antispam.HeuristicRule `yaml:"rules"`
}

type ListenerSettings struct {
	AcceptHost string `yaml:"socket_host"`
	AcceptPort string `yaml:"socket_port"`

	GithubWebhookEndpoint string `yaml:"github_webhook_endpoint"`
	GithubWebhookSecret   string `yaml:"github_webhook_secret"`
}

type Role struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
}
