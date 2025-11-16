// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package util

import (
	"bytes"
	"embed"
	"strings"
	"text/template"
)

func EscapeMarkdown(str string) string {
	return strings.NewReplacer(
		"`", "\\`",
		"_", "\\_",
		"*", "\\*",
		"~", "\\~",
		">", "\\>",
	).Replace(str)
}

func EmbedDescription(tplPath string,
	tplData embed.FS,
	fields map[string]any,
) string {
	var buf bytes.Buffer

	tpl, _ := template.ParseFS(tplData, tplPath)
	_ = tpl.Execute(&buf, fields)

	return buf.String()
}
