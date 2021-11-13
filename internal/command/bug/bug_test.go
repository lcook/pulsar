/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package command

import (
	"regexp"
	"testing"
)

func TestBugRegex(t *testing.T) {
	var tt = []struct {
		str      string
		expected int
	}{
		/*
		 * Cases that should not match.
		 */
		{"Lorem ipsum", 0},
		{"LOREm bug IPSUM !451", 0},
		{"HeLLo WOrLD BuG !1819", 0},
		{"hello world bug !!100", 0},
		{"bug !abc", 0},
		{"bug !AbC", 0},
		{"bug! !abc", 0},
		{"bug!19", 0},
		{"bug ! 391", 0},
		{"bug  !100", 0},
		/*
		 * Cases that should match.
		 */
		{"lorem ipsum bug !500!", 2},
		{"LORem IPSUm bug !194", 2},
		{"heLLo bug !414 WORLD", 2},
		{"hello world      bug !50", 2},
		{"bug !1 baz", 2},
		{"b u g bug !1019", 2},
		{"a aaa aaaaa bu   bug !  bug!! bug !849", 2},
	}
	reg := regexp.MustCompile(bugzRegex)
	for _, tc := range tt {
		actual := len(reg.FindStringSubmatch(tc.str))
		if actual != tc.expected {
			t.Errorf("expected %d, got %d", tc.expected, actual)
		}
	}

}
