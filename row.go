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
