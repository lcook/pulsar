/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package command

const (
	prefix string = "!"

	embedColorFreeBSD int = 0xEB0028
)

type Command struct {
	Name        string
	Description string
	Handler     any
}
