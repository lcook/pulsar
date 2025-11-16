// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package version

// Default version displayed with the (-v) pulsar version
// command-line flag. This is set during the build phase,
// usually including the current git commit short-hash.
var Build string = "devel"
