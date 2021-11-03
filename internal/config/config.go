/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port  string
	Token string
}

func Load(c string) (Config, error) {
	file, err := os.Open(c)
	if err != nil {
		return Config{}, err
	}
	//nolint
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)

	if err != nil {
		return Config{}, err
	}

	return Config{cfg.Port, cfg.Token}, nil
}
