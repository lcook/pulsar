/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package relay

import "net/http"

type Hook interface {
	Response(any) func(w http.ResponseWriter, r *http.Request)
	LoadConfig(string) error
	Endpoint() string
	Options() byte
}
