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

	"github.com/bsdlabs/pulsar/internal/util"
)

const (
	tplReportPath string = "templates/report.tpl"
)

//go:embed templates/report.tpl
var tplReportData embed.FS

func embedReport(bgg *bug) string {
	return util.EmbedDescription(tplReportPath, tplReportData, map[string]any{
		"status": bgg.Status,
		"product": func(str string) string {
			/*
			 * We cannot escape `&` in the Description field,
			 * so instead replace with `and`.  If not, it's
			 * character reference will be displayed `&amp;`.
			 */
			return strings.NewReplacer(
				"&", "and",
			).Replace(str)
		}(bgg.Product),
		"component": bgg.Component,
		"summary":   util.EscapeMarkdown(bgg.Summary),
		"url":       fmt.Sprintf(bugzReport, bgg.ID),
	})
}
