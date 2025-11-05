/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
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

	h.commands = []Command{
		{
			Name:        "help",
			Description: "Show this help page",
			Handler:     h.Help,
		},
		{
			Name:        "role",
			Description: "Assign yourself to a defined role",
			Handler:     h.Role,
		},
		{
			Name:        "bug",
			Description: "Display information of a Bugzilla report providing an ID",
			Handler:     h.Bug,
		},
		{
			Name:        "status",
			Description: "Display bot status",
			Handler:     h.Status,
		},
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
