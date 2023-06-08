/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package relay

import "net/http"

type Hook interface {
	Response(any) func(w http.ResponseWriter, r *http.Request)
	LoadConfig(string) error
	Endpoint() string
	Options() byte
}
