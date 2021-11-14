package util

import (
	"bytes"
	"embed"
	"strings"
	"text/template"
)

func EscapeMarkdown(s string) string {
	return strings.NewReplacer(
		"`", "\\`",
		"_", "\\_",
		"*", "\\*",
		"~", "\\~",
		">", "\\>",
	).Replace(s)
}

func EmbedDescription(tplPath string, tplData embed.FS,
	fields map[string]interface{}) string {
	tpl, _ := template.ParseFS(tplData, tplPath)
	var buf bytes.Buffer
	_ = tpl.Execute(&buf, fields)
	return buf.String()
}
