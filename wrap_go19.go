// +build go1.9

package sqle

import (
	"database/sql"
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
	case *sql.Conn:
		db.DB = *(**sql.DB)(unsafe.Pointer(obj))
		return &Conn{Conn: obj, db: db}, nil
	}
	return nil, nil
}
