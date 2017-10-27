package sqle

import (
	"database/sql"
	"reflect"
	"unsafe"
)

func Wrap(object interface{}, options ...DBOption) (_ interface{}, err error) {
	db := &DB{dbOptions: &dbOptions{mapper: defaultMapper}}
	for _, option := range options {
		if err = option(db.dbOptions); err != nil {
			return nil, err
		}
	}
	switch obj := object.(type) {
	case *sql.DB:
		db.DB = obj
		return db, nil
	case *sql.Row:
		rows := *(**sql.Rows)(unsafe.Pointer(uintptr(unsafe.Pointer(obj)) + offsetRowRows))
		db.DB = **(***sql.DB)(unsafe.Pointer(rows))
		return &Row{rows: &Rows{Rows: rows, db: db}}, nil
	case *sql.Rows:
		db.DB = **(***sql.DB)(unsafe.Pointer(obj))
		return &Rows{Rows: obj, db: db}, nil
	case *sql.Tx:
		db.DB = *(**sql.DB)(unsafe.Pointer(obj))
		return &Tx{Tx: obj, db: db}, nil
	case *sql.Stmt:
		db.DB = *(**sql.DB)(unsafe.Pointer(obj))
		return &Stmt{Stmt: obj, db: db}, nil
	}
	return nil, nil
}

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
