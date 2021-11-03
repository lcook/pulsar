/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
	"time"

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

func (c *commit) embed(repo, branch string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: c.embedDescription(repo, branch),
		Color:       color,
	}
}

//go:embed templates/commit.tpl
var tplData embed.FS

func (c *commit) embedDescription(repo, branch string) string {
	tpl, _ := template.ParseFS(tplData, "templates/commit.tpl")
	var msg bytes.Buffer
	_ = tpl.Execute(&msg, map[string]interface{}{
		"reponame":   repo,
		"gitrepo":    c.gitRepo(repo),
		"branchname": branch,
		"gitbranch":  c.gitBranch(repo, branch),
		"summary": func(s string) string {
			shortlog := strings.Split(s, "\n")[0]
			markdown := strings.NewReplacer(
				"`", "\\`",
				"_", "\\_",
				"*", "\\*",
				"~", "\\~",
			)
			return markdown.Replace(shortlog)
		}(c.Message),
		"committer": c.Committer.Name,
		"hash":      c.shortHash(),
		"gitcommit": c.gitCommit(repo),
	})
	return msg.String()
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
