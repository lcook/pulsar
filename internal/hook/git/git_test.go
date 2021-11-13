/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import (
	"fmt"
	"testing"
)

var (
	gitCommit      string = "12a61f4e173fb3a11c05d64"
	gitCommitShort string = gitCommit[0:7]
)

func TestGitRepo(t *testing.T) {
	var tt = []struct {
		commit   commit
		repo     string
		expected string
	}{
		{commit{}, "ports", fmt.Sprintf(cgitRepo, "ports")},
		{commit{}, "src", fmt.Sprintf(cgitRepo, "src")},
		{commit{}, "docs", fmt.Sprintf(cgitRepo, "docs")},
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
		{commit{}, "ports", "2021Q4", fmt.Sprintf(cgitBranch, "ports", "2021Q4")},
		{commit{}, "src", "stable/13", fmt.Sprintf(cgitBranch, "src", "stable/13")},
		{commit{}, "docs", "main", fmt.Sprintf(cgitBranch, "docs", "main")},
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
		{commit{ID: gitCommit}, "ports", fmt.Sprintf(cgitCommit, "ports", gitCommit)},
		{commit{ID: gitCommit}, "src", fmt.Sprintf(cgitCommit, "src", gitCommit)},
		{commit{ID: gitCommit}, "docs", fmt.Sprintf(cgitCommit, "docs", gitCommit)},
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
		{commit{ID: gitCommit}, gitCommitShort},
	}
	for _, tc := range tt {
		actual := tc.commit.shortHash()
		if actual != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, actual)
		}
	}
}
