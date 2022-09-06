/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2022, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
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
