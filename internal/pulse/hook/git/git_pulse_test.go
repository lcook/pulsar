/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package git

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidHmac(t *testing.T) {
	var (
		secret  string = "deadbeef"
		payload []byte = []byte("webhook data")
		tt             = struct {
			p        Pulse
			expected bool
		}{
			Pulse{}, true,
		}
	)
	tt.p.Config.Secret = secret
	hm := hmac.New(sha1.New, []byte(secret))
	hm.Write(payload)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Add("X-Hub-Signature", "sha1="+hex.EncodeToString(hm.Sum(nil)))
	w := httptest.NewRecorder()
	if tt.p.validHmac(payload, w, req) != tt.expected {
		t.Error()
	}
}
