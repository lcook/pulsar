/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package relay

/*
 * Optional middleware a hook can use for convenience.
 */
const (
	/*
	 * Check whether the incoming method is `POST`.
	 */
	OptionCheckMethod byte = 1 << iota
	/*
	 * Check whether the application type sent is `application/json`.
	 */
	OptionCheckType
	/*
	 * Reasonable defaults for webhook listening.
	 */
	DefaultOptions = (OptionCheckMethod | OptionCheckType)
)
