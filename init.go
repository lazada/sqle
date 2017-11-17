package sqle

import (
	"database/sql"
	"reflect"
)

func init() {
	var (
		rowtyp = reflect.TypeOf(sql.Row{})
		found  bool
	)
	for i, n := 0, rowtyp.NumField(); i < n; i++ {
		if field := rowtyp.Field(i); field.Name == `rows` {
			offsetRowRows, found = field.Offset, true
		}
	}
	if !found {
		panic(`unexpected structure of database/sql/Row`)
	}
}

var (
	offsetRowRows uintptr
)
