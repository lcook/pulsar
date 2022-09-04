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

func EmbedDescription(tplPath string, tplData embed.FS,
	fields map[string]any) string {
	tpl, _ := template.ParseFS(tplData, tplPath)

	var buf bytes.Buffer

	_ = tpl.Execute(&buf, fields)

	return buf.String()
}
