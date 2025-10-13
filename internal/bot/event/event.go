/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package event

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	embedDeleteColor int = 0xDC322F
	embedUpdateColor int = 0x268BD2

	maxContentLength int    = 1024
	maxContentMarker string = "\n\n<truncated>"

	eventLogChannel string = ""
)

func truncateContent(content string) string {
	if len(content) > maxContentLength {
		return content[:maxContentLength-len(maxContentMarker)] + maxContentMarker
	}

	return content
}

func buildContentField(content string, attachments []*discordgo.MessageAttachment) string {
	var builder strings.Builder
	if content != "" {
		builder.WriteString(content)
	}

	for idx, attachment := range attachments {
		if idx != 0 || builder.Len() > 0 {
			builder.WriteByte('\n')
		}

		builder.WriteString(fmt.Sprintf("<%s (%s)>", attachment.Filename, attachment.ContentType))
	}

	return truncateContent(builder.String())
}

func auditLogActions(session *discordgo.Session, member *discordgo.Member, action discordgo.AuditLogAction) ([]*discordgo.AuditLogEntry, error) {
	log, err := session.GuildAuditLog(member.GuildID, "", "", int(action), 100)
	if err != nil {
		return nil, err
	}

	var entries []*discordgo.AuditLogEntry

	for _, entry := range log.AuditLogEntries {
		if entry.TargetID == member.User.ID {
			entries = append(entries, entry)
		}
	}

	return entries, err
}
