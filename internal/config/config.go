/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	BotSettings      `yaml:"bot"`
	ListenerSettings `yaml:"listener"`
}

func FromFile[T any](path string) (T, error) {
	var settings T

	contents, err := os.ReadFile(path)
	if err != nil {
		return settings, err
	}

	err = yaml.Unmarshal(contents, &settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}
