// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package antispam

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

type TestCase struct {
	Name     string `json:"name"`
	Start    string `json:"start"`
	Messages []struct {
		Content   string `json:"content"`
		ChannelID string `json:"channel_id"`
		Timestamp string `json:"timestamp"`
	} `json:"messages"`
	Match string `json:"match"`
}

func TestEvaluateRules(t *testing.T) {
	definitions := LoadTestdata[[]HeuristicRule](t, "definitions.yaml")

	for _, tcase := range LoadTestdata[[]TestCase](t, "cases.json") {
		t.Run(tcase.Name, func(t *testing.T) {
			logs := make([]*Log, 0, len(tcase.Messages))

			for _, message := range tcase.Messages {
				ts, err := time.Parse(time.RFC3339, message.Timestamp)
				if err != nil {
					t.Fatalf(
						"invalid timestamp %q in testcase %q: %v",
						message.Timestamp,
						tcase.Name,
						err,
					)
				}

				logs = append(logs, &Log{
					Message: &discordgo.Message{
						ChannelID: message.ChannelID,
						Content:   message.Content,
						Timestamp: ts,
					},
					Hash: func() string {
						sha := sha512.New()
						sha.Write([]byte(message.Content))

						return hex.EncodeToString(sha.Sum(nil))
					}(),
				})
			}

			var rule HeuristicRule

			last := logs[len(logs)-1]
			for _, definition := range definitions {
				results := evaluateRule(
					last.Hash,
					logs,
					definition,
					last.Message.Timestamp,
				)
				if len(results) > 0 {
					rule = definition
					break
				}
			}

			if rule.ID != tcase.Match {
				t.Fatalf(
					"case %q: expected match=%v got match=%v",
					tcase.Name,
					tcase.Match,
					rule.ID,
				)
			}
		})
	}
}

func LoadTestdata[T any](t *testing.T, file string) T {
	t.Helper()

	path := filepath.Join("testdata", file)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read testdata file %s: %v", path, err)
	}

	var (
		result    T
		unmarshal func([]byte, any) error
	)

	switch filepath.Ext(file) {
	case ".yaml", ".yml":
		unmarshal = yaml.Unmarshal
	case ".json":
		unmarshal = json.Unmarshal
	default:
		t.Fatalf("unsupported file tyle: %s", file)
	}

	if err := unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal testdata file %s: %v", path, err)
	}

	return result
}
