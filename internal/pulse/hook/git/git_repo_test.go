/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import (
	"testing"
)

func TestCleanRepo(t *testing.T) {
	var tt = []struct {
		repo     repository
		expected string
	}{
		{repository{Name: "freebsd-ports"}, "ports"},
		{repository{Name: "freebsd-src"}, "src"},
		{repository{Name: "freebsd-docs"}, "docs"},
	}
	for _, tc := range tt {
		if tc.repo.String() != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, tc.repo.String())
		}
	}
}
