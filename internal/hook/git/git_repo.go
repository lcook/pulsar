/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import (
	"strings"
)

type repository struct {
	ID          int    `json:"id,omitempty"`
	NodeID      string `json:"node_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Fullname    string `json:"full_name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (r *repository) cleanRepo() {
	/*
	* Arbitrarily trim the repository name prefix
	* `freebsd-` since we track the FreeBSD GitHub
	* mirror repositories.
	*
	* `freebsd-ports`, `freebsd-src` and `freebsd-docs`
	* are returned as `ports`, `src`, and `docs`
	* respectively.
	 */
	r.Name = strings.TrimPrefix(r.Name, "freebsd-")
}
