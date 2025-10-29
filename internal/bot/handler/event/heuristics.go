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
		Messages int           `toml:"messages"`
		Channels int           `toml:"channels"`
		Mentions int           `toml:"mentions"`
		Window   time.Duration `toml:"Window"`
	} `toml:"thresholds`
	Timeout time.Duration `toml:"timeout"`
}

type Heuristics struct {
	Rules []HeuristicRule `yaml:"rules"`
}

//go:embed data/heuristics.yaml
var heuristicsData []byte

func GetHeuristics(hash string, logs []Log, rules []HeuristicRule) ([]Log, *HeuristicRule) {
	var rule *HeuristicRule

	var spamLogs []Log

	now := time.Now().UTC()

	for _, r := range rules {
		var target []Log

		for _, log := range logs {
			if now.Sub(log.Message.Timestamp.UTC()) > r.Thresholds.Window {
				continue
			}

			target = append(target, log)
		}

		if r.Duplicated {
			var filtered []Log
			for _, log := range target {
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
			for _, log := range target {
				channels[log.Message.ChannelID] = struct{}{}
			}

			if len(channels) < r.Thresholds.Channels {
				continue
			}
		}

		if r.Thresholds.Mentions > 0 {
			var matched bool

			re := regexp.MustCompile(`<@!?(\d+)>`)
			for _, log := range target {
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
