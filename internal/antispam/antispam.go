/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package antispam

import (
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
	"github.com/lcook/pulsar/internal/cache"
)

type Log struct {
	Message *discordgo.Message
	Hash    string

	deleted atomic.Bool
}

func (l *Log) Deleted() bool { return l.deleted.Load() }
func (l *Log) MarkDeleted()  { l.deleted.Store(true) }

func Run(m *discordgo.MessageCreate, hash string, cache *cache.RingBuffer[Log], rules []HeuristicRule) ([]*Log, *HeuristicRule) {
	var logs []*Log

	for idx := range cache.Slice() {
		log := &cache.Slice()[idx]
		if m.Author.ID == log.Message.Author.ID && !log.Deleted() {
			logs = append(logs, log)
		}
	}

	spamLogs, rule := evaluateRules(hash, logs, rules)
	if len(spamLogs) == 0 || rule == nil {
		return spamLogs, rule
	}

	return spamLogs, rule
}
