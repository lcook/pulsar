/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package event

import (
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
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
	Logs     *cache.RingBuffer[Log]
}

type Log struct {
	Message *discordgo.Message
	Hash    string

	deleted atomic.Bool
}

func New(settings config.Settings, buffer uint64) *Handler {
	h := &Handler{
		Settings: settings,
		Logs:     cache.NewRingBuffer[Log](buffer),
	}

	h.Events = append(h.Events, h.MessageCreate)
	h.Events = append(h.Events, h.MessageDelete)
	h.Events = append(h.Events, h.MessageUpdate)
	h.Events = append(h.Events, h.GuildMemberRemove)
	h.Events = append(h.Events, h.AutoModExecution)

	return h
}
