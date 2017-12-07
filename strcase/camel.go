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
	"unicode"
	"unicode/utf8"
)

func ToCamel(str string) string {
	if str == `` {
		return ``
	}
	buf := make([]byte, len(str))
	pos, up := 0, true

	for _, r := range str {
		if r == '_' || r == '-' || unicode.IsSpace(r) {
			up = true
			continue
		}
		if up && unicode.IsLetter(r) {
			up, r = false, unicode.ToUpper(r)
		}
		if r >= utf8.RuneSelf {
			pos += utf8.EncodeRune(buf[pos:], r)
		} else {
			buf[pos], pos = byte(r), pos+1
		}
	}
	return string(buf[:pos])
}
