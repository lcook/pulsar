/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) Lewis Cook <lcook@FreeBSD.org>
 */
package command

import (
	"testing"
)

const (
	nomatch int = 0
	match   int = 2
)

func TestBugRegex(t *testing.T) {
	// var s = func(s string) string { return fmt.Sprintf(s, bugzRegexPrefix) }
	// var tt = []struct {
	// 	str      string
	// 	expected int
	// }{
	// 	/*
	// 	 * Cases that _should not_ match.
	// 	 */
	// 	{"Lorem ipsum", nomatch},
	// 	{s("LOREm bug IPSUM %s451"), nomatch},
	// 	{s("HeLLo WOrLD BuG %s1819"), nomatch},
	// 	{s("hello world bug %[1]s%[1]s100"), nomatch},
	// 	{s("bug %sabc"), nomatch},
	// 	{s("bug %sAbC"), nomatch},
	// 	{s("bug%[1]s %[1]sabc"), nomatch},
	// 	{s("bug%s19"), nomatch},
	// 	{s("bug %s 391"), nomatch},
	// 	{s("bug  %s100"), nomatch},
	// 	/*
	// 	 * Cases that _should_ match.
	// 	 */
	// 	{s("lorem ipsum bug %[1]s500%[1]s"), match},
	// 	{s("LORem IPSUm bug %s194"), match},
	// 	{s("heLLo bug %s414 WORLD"), match},
	// 	{s("hello world      bug %s50"), match},
	// 	{s("bug %s1 baz"), match},
	// 	{s("b u g bug %s1019"), match},
	// 	{s("a aaa aaaaa bu   bug %[1]s  bug%[1]s%[1]s bug %[1]s849"), match},
	// }
	// reg := regexp.MustCompile(bugzRegex)
	// for _, tc := range tt {
	// 	actual := len(reg.FindStringSubmatch(tc.str))
	// 	if actual != tc.expected {
	// 		t.Errorf("expected %d, got %d", tc.expected, actual)
	// 	}
	// }
}
