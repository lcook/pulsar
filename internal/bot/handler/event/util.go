// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package event

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	maxContentLength int    = 1024
	maxContentMarker string = "\n\n<truncated>"
)

func truncateContent(content string) string {
	if len(content) > maxContentLength {
		return content[:maxContentLength-len(maxContentMarker)] + maxContentMarker
	}

	return content
}

func hashContent(content string) string {
	sha := sha512.New()
	sha.Write([]byte(content))

	return hex.EncodeToString(sha.Sum(nil))
}

func buildContentField(
	content string,
	attachments []*discordgo.MessageAttachment,
) string {
	var builder strings.Builder
	if content != "" {
		builder.WriteString(content)
	}

	for idx, attachment := range attachments {
		if idx != 0 || builder.Len() > 0 {
			builder.WriteByte('\n')
		}

		builder.WriteString(
			fmt.Sprintf(
				"<%s (%s)>",
				attachment.Filename,
				attachment.ContentType,
			),
		)
	}

	return truncateContent(builder.String())
}

func auditLogActions(
	session *discordgo.Session,
	member *discordgo.Member,
	action discordgo.AuditLogAction,
	limit int,
) ([]*discordgo.AuditLogEntry, error) {
	log, err := session.GuildAuditLog(
		member.GuildID,
		"",
		"",
		int(action),
		limit,
	)
	if err != nil {
		return nil, err
	}

	var entries []*discordgo.AuditLogEntry

	for _, entry := range log.AuditLogEntries {
		if entry.TargetID == member.User.ID {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func auditLogActionsLast(
	session *discordgo.Session,
	member *discordgo.Member,
	action discordgo.AuditLogAction,
	limit int,
	window time.Duration,
) ([]*discordgo.AuditLogEntry, error) {
	var (
		err     error
		entries []*discordgo.AuditLogEntry
	)
	// Discord may lag a few hundred ms before the audit entry is written.
	// Retry up to three times with a truncated exponential back-off (100 -> 300 ms)
	// until the desired entry appears or we give up.
	for attempt := range 3 {
		entries, err = auditLogActions(session, member, action, limit)
		if err == nil && len(entries) > 0 {
			break
		}

		if attempt < 2 {
			time.Sleep(time.Duration(100*(attempt+1)) * time.Millisecond)
		}
	}

	if err != nil {
		return entries, err
	}

	cutoff := time.Now().UTC().Add(-window)

	logs := make([]*discordgo.AuditLogEntry, 0, len(entries))
	for _, entry := range entries {
		timestamp, err := discordgo.SnowflakeTimestamp(entry.ID)
		if err != nil || timestamp.IsZero() {
			continue
		}
		// Drop entries older than the supplied window.
		// Without this, a past kick/ban audit entry could be
		// re-used if the user re-joins and later leaves again.
		if timestamp.Before(cutoff) {
			continue
		}

		logs = append(logs, entry)
	}

	return logs, nil
}

func canViewChannel(
	session *discordgo.Session,
	guildID, channelID string,
) bool {
	everyone, _ := session.State.Role(guildID, guildID)
	channel, _ := session.Channel(channelID)

	for _, permission := range channel.PermissionOverwrites {
		if permission.ID == everyone.ID &&
			permission.Deny&discordgo.PermissionViewChannel != 0 {
			return false
		}
	}

	return true
}

func logMember(
	member *discordgo.User,
	level log.Level,
	message string,
	fields ...log.Fields,
) {
	created, _ := discordgo.SnowflakeTimestamp(member.ID)

	logFields := log.Fields{
		"id":       member.ID,
		"username": member.Username,
		"nickname": member.DisplayName(),
		"created":  created,
		"verified": member.Verified,
	}

	for _, _fields := range fields {
		maps.Copy(logFields, _fields)
	}

	logEntry := log.WithFields(logFields)

	switch level {
	case log.TraceLevel:
		logEntry.Trace(message)
	case log.DebugLevel:
		logEntry.Debug(message)
	case log.InfoLevel:
		logEntry.Info(message)
	case log.WarnLevel:
		logEntry.Warn(message)
	}
}
