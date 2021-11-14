/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package command

import (
	"embed"
	"fmt"
	"strings"

	"github.com/bsdlabs/pulseline/internal/util"
)

const (
	tplReportPath string = "templates/report.tpl"
)

//go:embed templates/report.tpl
var tplReportData embed.FS

func embedReport(b bug) string {
	return util.EmbedDescription(tplReportPath, tplReportData, map[string]interface{}{
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
		"summary":   util.EscapeMarkdown(b.Summary),
		"url":       fmt.Sprintf(bugzReport, b.ID),
	})
}
