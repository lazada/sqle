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
	"errors"
)

// Row is the result of calling QueryRow to select a single row.
type Row struct {
	err  error
	rows *Rows
}

// Scan copies the columns from the matched row into the values
// pointed at by dest. See the documentation on Rows.Scan for details.
// If more than one row matches the query, Scan uses the first row and discards the rest.
// If no row matches the query, Scan returns sql.ErrNoRows.
func (r *Row) Scan(dest ...interface{}) (err error) {
	if r.err != nil {
		return r.err
	}
	defer r.rows.Close()

	for _, d := range dest {
		if _, ok := d.(*sql.RawBytes); ok {
			return errors.New("sql: RawBytes isn't allowed on Row.Scan")
		}
	}
	if !r.rows.Next() {
		if err = r.rows.Err(); err != nil {
			return
		}
		return sql.ErrNoRows
	}
	if err = r.rows.Scan(dest...); err == nil {
		err = r.rows.Close()
	}
	return
}

var (
	ErrMiss = errors.New(`miss`)
)
