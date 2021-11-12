/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import "testing"

func TestGitRepo(t *testing.T) {
	var tt = []struct {
		commit   commit
		repo     string
		expected string
	}{
		{commit{}, "ports", "https://cgit.freebsd.org/ports/"},
		{commit{}, "src", "https://cgit.freebsd.org/src/"},
		{commit{}, "docs", "https://cgit.freebsd.org/docs/"},
	}
	for _, tc := range tt {
		actual := tc.commit.gitRepo(tc.repo)
		if actual != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, actual)
		}
	}
}

func TestGitBranch(t *testing.T) {
	var tt = []struct {
		commit   commit
		repo     string
		branch   string
		expected string
	}{
		{commit{}, "ports", "2021Q4", "https://cgit.freebsd.org/ports/?h=2021Q4"},
		{commit{}, "src", "stable/13", "https://cgit.freebsd.org/src/?h=stable/13"},
		{commit{}, "docs", "main", "https://cgit.freebsd.org/docs/?h=main"},
	}
	for _, tc := range tt {
		actual := tc.commit.gitBranch(tc.repo, tc.branch)
		if actual != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, actual)
		}
	}

}

func TestGitCommit(t *testing.T) {
	var tt = []struct {
		commit   commit
		repo     string
		expected string
	}{
		{commit{ID: "12a61f4e173fb3a11c05d64"}, "ports", "https://cgit.freebsd.org/ports/commit/?id=12a61f4e173fb3a11c05d64"},
		{commit{ID: "12a61f4e173fb3a11c05d64"}, "src", "https://cgit.freebsd.org/src/commit/?id=12a61f4e173fb3a11c05d64"},
		{commit{ID: "12a61f4e173fb3a11c05d64"}, "docs", "https://cgit.freebsd.org/docs/commit/?id=12a61f4e173fb3a11c05d64"},
	}
	for _, tc := range tt {
		actual := tc.commit.gitCommit(tc.repo)
		if actual != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, actual)
		}
	}
}

func TestShortHash(t *testing.T) {
	var tt = []struct {
		commit   commit
		expected string
	}{
		{commit{ID: "12a61f4e173fb3a11c05d64"}, "12a61f4"},
	}
	for _, tc := range tt {
		actual := tc.commit.shortHash()
		if actual != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, actual)
		}
	}
}
