/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import "testing"

func TestCleanRef(t *testing.T) {
	var tt = []struct {
		events   commitEvent
		expected string
	}{
		{commitEvent{Ref: "refs/heads/main"}, "main"},
		{commitEvent{Ref: "refs/heads/2021Q4"}, "2021Q4"},
		{commitEvent{Ref: "refs/heads/stable/13"}, "stable/13"},
	}
	for _, tc := range tt {
		tc.events.cleanRef()
		if tc.events.Ref != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, tc.events.Ref)
		}
	}

}
