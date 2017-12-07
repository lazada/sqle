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

package sqle

import (
	"database/sql"
	"reflect"
	"time"
)

// Rows is the result of a query. Its cursor starts before the first row of the result set.
type Rows struct {
	*sql.Rows
	db       *DB
	err      error
	columns  []string
	coltypes []*sql.ColumnType
	pointers []interface{}
}

// Columns returns the column names.
// Columns returns an error if the rows are closed, or if the rows
// are from QueryRow and there was a deferred error.
func (r *Rows) Columns() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.columns == nil {
		if r.columns, r.err = r.Rows.Columns(); r.err != nil {
			return nil, r.err
		}
		if r.pointers == nil {
			r.pointers = make([]interface{}, 0, len(r.columns))
		}
	}
	return r.columns, nil
}

// ColumnTypes returns column information such as column type, length, and nullable.
// Some information may not be available from some drivers.
func (r *Rows) ColumnTypes() ([]*sql.ColumnType, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.coltypes == nil {
		if r.coltypes, r.err = r.Rows.ColumnTypes(); r.err != nil {
			return nil, r.err
		}
		if r.pointers == nil {
			r.pointers = make([]interface{}, 0, len(r.coltypes))
		}
	}
	return r.coltypes, nil
}

// Scan copies the columns in the current row into the values pointed at by dest.
func (r *Rows) Scan(dest ...interface{}) (err error) {
	if len(dest) == 0 {
		return nil
	}
	if r.Columns(); r.err != nil {
		return r.err
	}
	var (
		mp       *map[string]interface{}
		a        []string
		c, m     int
		coltyp   reflect.Type
		colnum   = len(r.columns)
		destnum  = len(dest)
		destlast = destnum - 1
		ptrs     = r.pointers[:0]
	)
loop:
	for i, j := 0, 0; i < destnum && j < colnum; i++ {
		switch t := dest[i].(type) {
		case sql.Scanner:
			ptrs, j = append(ptrs, dest[i]), j+1
		case Pointers:
			if m, j = j, j+t.Num(); j > colnum || i == destlast {
				if r.db.strict {
					return ErrMiss
				}
				j = colnum
			}
			ptrs, m = t.Pointers(ptrs, r.columns[m:j])
			if r.db.strict && m > 0 {
				return ErrMiss
			}
		case map[string]interface{}:
			if r.ColumnTypes(); r.err != nil {
				return r.err
			}
			for c, m, mp = j, j, &t; c < colnum; c++ {
				if coltyp = r.coltypes[c].ScanType(); coltyp == nil {
					ptrs = append(ptrs, new(interface{}))
				} else {
					ptrs = append(ptrs, reflect.New(coltyp).Interface())
				}
			}
			break loop
		case *map[string]interface{}:
			if r.ColumnTypes(); r.err != nil {
				return r.err
			}
			for c, m, mp = j, j, t; c < colnum; c++ {
				if coltyp = r.coltypes[c].ScanType(); coltyp == nil {
					ptrs = append(ptrs, new(interface{}))
				} else {
					ptrs = append(ptrs, reflect.New(coltyp).Interface())
				}
			}
			break loop
		default:
			if a, err = r.db.mapper.Aliases(dest[i]); err == nil {
				if m, j = j, j+len(a); j > colnum || i == destlast {
					if r.db.strict {
						return ErrMiss
					}
					j = colnum
				}
				ptrs, m, err = r.db.mapper.Pointers(dest[i], ptrs, r.columns[m:j])
				if err == nil && r.db.strict && m > 0 {
					return ErrMiss
				}
			} else if err == ErrSrcNotPtr {
				return
			} else {
				ptrs, j = append(ptrs, dest[i]), j+1
			}
			err = nil
		}
	}
	if err = r.Rows.Scan(ptrs...); err == nil && mp != nil {
		for i, n := m, len(ptrs); i < n; i++ {
			switch ptrs[i].(type) {
			case *bool:
				(*mp)[r.columns[i]] = *(ptrs[i].(*bool))
			case *float32:
				(*mp)[r.columns[i]] = *(ptrs[i].(*float32))
			case *float64:
				(*mp)[r.columns[i]] = *(ptrs[i].(*float64))
			case *int8:
				(*mp)[r.columns[i]] = *(ptrs[i].(*int8))
			case *int16:
				(*mp)[r.columns[i]] = *(ptrs[i].(*int16))
			case *int32:
				(*mp)[r.columns[i]] = *(ptrs[i].(*int32))
			case *int64:
				(*mp)[r.columns[i]] = *(ptrs[i].(*int64))
			case *uint8:
				(*mp)[r.columns[i]] = *(ptrs[i].(*uint8))
			case *uint16:
				(*mp)[r.columns[i]] = *(ptrs[i].(*uint16))
			case *uint32:
				(*mp)[r.columns[i]] = *(ptrs[i].(*uint32))
			case *uint64:
				(*mp)[r.columns[i]] = *(ptrs[i].(*uint64))
			case *string:
				(*mp)[r.columns[i]] = *(ptrs[i].(*string))
			case *interface{}:
				(*mp)[r.columns[i]] = *(ptrs[i].(*interface{}))
			case *time.Time:
				(*mp)[r.columns[i]] = *(ptrs[i].(*time.Time))
			default:
				(*mp)[r.columns[i]] = reflect.ValueOf(ptrs[i]).Elem().Interface()
			}
		}
	}
	return
}
