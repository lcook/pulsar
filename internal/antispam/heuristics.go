// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package antispam

import (
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

func evaluate(
	hash string,
	logs []*Log,
	rules []HeuristicRule,
) ([]*Log, *HeuristicRule) {
	var (
		matchRule *HeuristicRule
		matchLogs []*Log
		timestamp = time.Now().UTC()
	)

	for _, rule := range rules {
		results := evaluateRule(hash, logs, rule, timestamp)
		if len(results) > 0 {
			matchLogs = results
			matchRule = &rule

			break
		}
	}

	return matchLogs, matchRule
}

func evaluateRule(
	hash string,
	logs []*Log,
	rule HeuristicRule,
	timestamp time.Time,
) []*Log {
	var target []*Log

	for idx := range logs {
		log := logs[idx]
		if timestamp.Sub(log.Message.Timestamp.UTC()) > rule.Thresholds.Window {
			continue
		}

		target = append(target, log)
	}

	if rule.Duplicated {
		var dupe []*Log

		for idx := range target {
			log := target[idx]
			if log.Hash == hash {
				dupe = append(dupe, log)
			}
		}

		target = dupe
	}

	if rule.Thresholds.Messages > 0 && len(target) < rule.Thresholds.Messages {
		return nil
	}

	if rule.Thresholds.Channels > 0 {
		channels := make(map[string]struct{})

		for idx := range target {
			log := target[idx]
			channels[log.Message.ChannelID] = struct{}{}
		}

		if len(channels) < rule.Thresholds.Channels {
			return nil
		}
	}

	if rule.Thresholds.Mentions > 0 {
		var matched bool

		re := regexp.MustCompile(`<@!?(\d+)>`)

		for idx := range target {
			log := target[idx]

			mentions := len(
				re.FindAllStringSubmatch(log.Message.Content, -1),
			)

			if mentions >= rule.Thresholds.Mentions {
				matched = true
				break
			}
		}

		if !matched {
			return nil
		}
	}

	return target
}
