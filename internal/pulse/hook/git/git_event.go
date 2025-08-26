/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package git

import (
	"encoding/json"
	"strings"
)

type commitEvent struct {
	Ref        string     `json:"ref,omitempty"`
	Before     string     `json:"before,omitempty"`
	After      string     `json:"after,omitempty"`
	Repository repository `json:"repository,omitempty"`
	Commits    []commit   `json:"commits,omitempty"`
}

func (ce *commitEvent) cleanRef() {
	/*
	 * Strim the raw reference prefix, leaving the branch
	 * name intact, e.g., main.
	 */
	ce.Ref = strings.TrimPrefix(ce.Ref, "refs/heads/")
}

func commitEventPayload(buf []byte) (*commitEvent, error) {
	var payload commitEvent

	err := json.Unmarshal(buf, &payload)
	if err != nil {
		return &commitEvent{}, err
	}

	payload.cleanRef()

	return &payload, nil
}
