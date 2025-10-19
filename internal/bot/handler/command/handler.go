/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package command

import (
	"github.com/lcook/pulsar/internal/config"
)

const (
	embedColorFreeBSD int = 0xEB0028
)

type Handler struct {
	Settings config.Settings

	commands []Command
}

type Command struct {
	Name        string
	Description string
	Handler     any
}

func New(settings config.Settings) *Handler {
	h := &Handler{Settings: settings}

	h.commands = []Command{
		{
			Name:        "help",
			Description: "Show help page",
			Handler:     h.Help,
		},
		{
			Name:        "role",
			Description: "Assign or remove a role to yourself",
			Handler:     h.Role,
		},
		{
			Name:        "bug",
			Description: "Fetch information with provided Bugzilla report ID",
			Handler:     h.Bug,
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
