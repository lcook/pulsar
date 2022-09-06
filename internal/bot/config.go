/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
 * All rights reserved.
 */
package bot

type Config struct {
	Port  string `json:"event_listener_port"`
	Token string `json:"discord_bot_token"`
}
