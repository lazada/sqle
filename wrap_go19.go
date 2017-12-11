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

// +build go1.9

package sqle

import (
	"database/sql"
	"unsafe"
)

// Wrap converts objects from the standard `database/sql` to the coresponding `sqle` objects.
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
