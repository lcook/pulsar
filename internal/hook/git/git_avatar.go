/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import (
	//nolint
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
)

const (
	gravatarBase      string = "https://www.gravatar.com/"
	gravatarIcon      string = gravatarBase + "/avatar/"
	gravatarIdenticon string = "?d=identicon"

	githubBase string = "https://github.com"
)

func Avatar(username, email string) string {
	avatar := fmt.Sprintf(githubBase+"/%s.png", username)
	//nolint
	if resp, _ := http.Get(avatar); resp.StatusCode != 200 {
		hash := md5.Sum([]byte(email))
		avatar = fmt.Sprintf(gravatarIcon+"%s.jpg%s", hex.EncodeToString(hash[:]), gravatarIdenticon)
		//nolint
		defer resp.Body.Close()
	}

	return avatar
}
