/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package event

import (
	_ "embed"
	"regexp"
	"time"
)

type HeuristicRule struct {
	ID         string `yaml:"id"`
	Duplicated bool   `yaml:"duplicated"`
	Thresholds struct {
		Messages int           `yaml:"messages"`
		Channels int           `yaml:"channels"`
		Mentions int           `yaml:"mentions"`
		Window   time.Duration `yaml:"window"`
	} `yaml:"thresholds"`
	Timeout time.Duration `yaml:"timeout"`
}

type Heuristics struct {
	Rules []HeuristicRule `yaml:"rules"`
}

//go:embed data/heuristics.yaml
var heuristicsData []byte

func GetHeuristics(hash string, logs []*Log, rules []HeuristicRule) ([]*Log, *HeuristicRule) {
	var rule *HeuristicRule

	var spamLogs []*Log

	now := time.Now().UTC()

	for _, r := range rules {
		var target []*Log

		for idx := range logs {
			log := logs[idx]
			if now.Sub(log.Message.Timestamp.UTC()) > r.Thresholds.Window {
				continue
			}

			target = append(target, log)
		}

		if r.Duplicated {
			var filtered []*Log

			for idx := range target {
				log := target[idx]
				if log.Hash == hash {
					filtered = append(filtered, log)
				}
			}

			target = filtered
		}

		if r.Thresholds.Messages > 0 && len(target) < r.Thresholds.Messages {
			continue
		}

		if r.Thresholds.Channels > 0 {
			channels := make(map[string]struct{})

			for idx := range target {
				log := target[idx]
				channels[log.Message.ChannelID] = struct{}{}
			}

			if len(channels) < r.Thresholds.Channels {
				continue
			}
		}

		if r.Thresholds.Mentions > 0 {
			var matched bool

			re := regexp.MustCompile(`<@!?(\d+)>`)

			for idx := range target {
				log := target[idx]

				mentions := len(re.FindAllStringSubmatch(log.Message.Content, -1))
				if mentions >= r.Thresholds.Mentions {
					matched = true
					break
				}
			}

			if !matched {
				continue
			}
		}

		rule = &r
		spamLogs = target

		break
	}

	return spamLogs, rule
}
