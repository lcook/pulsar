/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021-2022, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import (
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/bsdlabs/pulseline/internal/util"
	"github.com/bwmarrin/discordgo"
)

const (
	cgitBase   string = "https://cgit.freebsd.org"
	cgitRepo   string = cgitBase + "/%s/"
	cgitBranch string = cgitBase + "/%s/?h=%s"
	cgitCommit string = cgitBase + "/%s/commit/?id=%s"
)

type commit struct {
	ID        string    `json:"id,omitempty"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Author    author    `json:"author,omitempty"`
	Committer committer `json:"committer,omitempty"`
	Added     []string  `json:"added,omitempty"`
	Removed   []string  `json:"removed,omitempty"`
	Modified  []string  `json:"modified,omitempty"`
}

type author struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

type committer struct {
	author
}

func (c *committer) String() string { return c.Name }

func (c *commit) embed(repo, branch string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       color,
		Description: c.embedCommit(repo, branch),
		Author: &discordgo.MessageEmbedAuthor{
			Name:    c.Committer.String(),
			IconURL: c.Committer.Avatar(),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%s repository", repo),
		},
		Timestamp: c.Timestamp.Format(time.RFC3339),
	}
}

const (
	tplCommitPath string = "templates/commit.tpl"
)

//go:embed templates/commit.tpl
var tplCommitData embed.FS

func (c *commit) embedCommit(repo, branch string) string {
	return util.EmbedDescription(tplCommitPath, tplCommitData, map[string]any{
		"reponame":   repo,
		"gitrepo":    c.gitRepo(repo),
		"branchname": branch,
		"gitbranch":  c.gitBranch(repo, branch),
		"summary":    util.EscapeMarkdown(strings.Split(c.Message, "\n")[0]),
		"committer":  c.Committer.String(),
		"hash":       c.shortHash(),
		"gitcommit":  c.gitCommit(repo),
	})
}

func (c *commit) gitRepo(repo string) string {
	return fmt.Sprintf(cgitRepo, repo)
}

func (c *commit) gitBranch(repo, branch string) string {
	return fmt.Sprintf(cgitBranch, repo, branch)
}

func (c *commit) gitCommit(repo string) string {
	return fmt.Sprintf(cgitCommit, repo, c.ID)
}

func (c *commit) shortHash() string {
	return c.ID[0:7]
}
