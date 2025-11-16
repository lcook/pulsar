// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"github.com/lcook/pulsar/internal/antispam"
	"github.com/lcook/pulsar/internal/cache"
	"github.com/lcook/pulsar/internal/config"
)

const (
	embedDeleteColor int = 0xDC322F
	embedUpdateColor int = 0x268BD2
)

type Handler struct {
	Settings config.Settings
	Events   []any
	Logs     *cache.RingBuffer[antispam.Log]
}

func New(settings config.Settings, buffer uint64) *Handler {
	h := &Handler{
		Settings: settings,
		Logs:     cache.NewRingBuffer[antispam.Log](buffer),
	}

	h.Events = append(h.Events, h.MessageCreate)
	h.Events = append(h.Events, h.MessageDelete)
	h.Events = append(h.Events, h.MessageUpdate)
	h.Events = append(h.Events, h.GuildMemberAdd)
	h.Events = append(h.Events, h.GuildMemberRemove)
	h.Events = append(h.Events, h.AutoModExecution)

	return h
}
