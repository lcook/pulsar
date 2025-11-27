// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package command

import (
	"time"

	"github.com/lcook/pulsar/internal/config"
)

const (
	embedColorFreeBSD int = 0xEB0028
)

type Handler struct {
	Settings config.Settings
	Started  time.Time

	commands []Command
}

type Command struct {
	Name        string
	Description string
	Handler     any
}

func New(settings config.Settings) *Handler {
	h := &Handler{
		Settings: settings,
		Started:  time.Now(),
	}

	available := map[string]Command{
		"help": {"help", "Show this help page", h.Help},
		"role": {"role", "Assign yourself to a defined role", h.Role},
		"bug": {
			"bug",
			"Display information of a Bugzilla report providing an ID",
			h.Bug,
		},
		"status": {"status", "Display bot status", h.Status},
		"review": {
			"review",
			"Display information of a Phabricator differential revision providing an ID",
			h.Review,
		},
	}

	for _, name := range settings.Commands {
		if cmd, ok := available[name]; ok {
			h.commands = append(h.commands, cmd)
		}
	}

	return h
}

func (h *Handler) Handlers() []any {
	hnd := make([]any, 0, len(h.commands))
	for _, handler := range h.commands {
		hnd = append(hnd, handler.Handler)
	}

	return hnd
}
