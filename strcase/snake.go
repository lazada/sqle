// Copyright 2017 Lazada South East Asia Pte. Ltd.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strcase

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func ToKebab(str string) string { return Snake(str, '-', 0) }

func ToSnake(str string) string { return Snake(str, '_', 0) }

func Snake(str string, sep rune, greed uint8) string {
	if str == `` {
		return ``
	}
	var (
		state, g  uint8
		pos, mark int
		r         rune
		s         = utf8.RuneLen(sep)
		buf       = make([]byte, len(str)+(len(str)>>1)*s)
	)
	for _, r = range str {
		switch {
		case r == sep || r == '_' || r == '-' || unicode.IsSpace(r):
			if state == 0 && g >= greed {
				continue
			}
			state, g, r = 0, g+1, sep
		case unicode.IsUpper(r):
			switch state {
			case 1:
				mark = pos
			case 2:
				if s > 1 {
					pos += utf8.EncodeRune(buf[pos:], sep)
				} else {
					buf[pos], pos = byte(sep), pos+1
				}
			}
			state, g = 1, 0
		case unicode.IsLower(r):
			if mark > 0 {
				copy(buf[mark+s:], buf[mark:])
				if r >= utf8.RuneSelf {
					pos += utf8.EncodeRune(buf[mark:], sep)
				} else {
					buf[mark], pos = byte(sep), pos+1
				}
			}
			state, g, mark = 2, 0, 0
		}
		if r >= utf8.RuneSelf {
			pos += utf8.EncodeRune(buf[pos:], r)
		} else {
			buf[pos], pos = byte(r), pos+1
		}
	}
	if greed == 0 && (r == sep || r == '_' || r == '-' || unicode.IsSpace(r)) {
		pos -= s
	}
	return strings.ToLower(string(buf[:pos]))
}
