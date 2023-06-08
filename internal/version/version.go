/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package version

var (
	/*
	 * Default version displayed in the (-v) pulsar version
	 * command-line flag.  This is set during the build phase,
	 * possibly including the current git commit short-hash.
	 */
	Build string = "devel"
)
