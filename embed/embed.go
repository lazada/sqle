package embed

import (
	"errors"
)

type DummyField struct{}

func (f *DummyField) Scan(value interface{}) error { return nil }

var (
	ErrValueNotFound = errors.New(`sqle: value not found`)
)
