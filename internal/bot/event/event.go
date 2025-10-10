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

	builderStr := builder.String()
	if len(builderStr) > maxContentLength {
		builderStr = builderStr[:maxContentLength-len(maxContentMarker)] + maxContentMarker
	}

	return builderStr
}
