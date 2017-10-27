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
