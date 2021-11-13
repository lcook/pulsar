/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package command

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed templates/report.tpl
var tplData embed.FS

func embedDescription(b bug) string {
	tpl, _ := template.ParseFS(tplData, "templates/report.tpl")
	var msg bytes.Buffer
	_ = tpl.Execute(&msg, map[string]interface{}{
		"status": b.Status,
		"product": func(s string) string {
			/*
			 * We cannot escape `&` in the Description field,
			 * so instead replace with `and`.  If not, it's
			 * character reference will be displayed `&amp;`.
			 */
			return strings.NewReplacer(
				"&", "and",
			).Replace(s)
		}(b.Product),
		"component": b.Component,
		"summary": func(s string) string {
			markdown := strings.NewReplacer(
				"`", "\\`",
				"_", "\\_",
				"*", "\\*",
				"~", "\\~",
			)
			return markdown.Replace(s)
		}(b.Summary),
		"url": fmt.Sprintf(bugzReport, b.ID),
	})
	return msg.String()
}
