// +build !debug

package sqle

import "testing"

func debug(t *testing.T, args ...interface{}) {}

func debugf(t *testing.T, fmt string, args ...interface{}) {}
