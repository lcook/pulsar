/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Port  string `json:"event_listener_port"`
	Token string `json:"discord_bot_token"`
}

func Load(c string) (Config, error) {
	file, err := os.Open(c)
	if err != nil {
		return Config{}, err
	}
	//nolint
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)

	if err != nil {
		return Config{}, err
	}

	return Config{cfg.Port, cfg.Token}, nil
}
