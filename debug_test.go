// +build debug

package sqle

import "testing"

func debug(t *testing.T, args ...interface{}) {
	t.Log(args...)
}

func debugf(t *testing.T, fmt string, args ...interface{}) {
	t.Logf(fmt, args...)
}
