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
	"context"
	"database/sql"
)

// Stmt is a prepared statement.
// A Stmt is safe for concurrent use by multiple goroutines.
type Stmt struct {
	*sql.Stmt
	db *DB
}

// Query executes a prepared query statement with given arguments
// and returns the query results as *Rows.
func (s *Stmt) Query(args ...interface{}) (*Rows, error) {
	return s.QueryContext(context.Background(), args...)
}

// QueryContext executes a prepared query statement with given arguments
// and returns the query results as *Rows.
func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*Rows, error) {
	rows, err := s.Stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: rows, db: s.db}, nil
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return sql.ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards the rest.
func (s *Stmt) QueryRow(args ...interface{}) *Row {
	return s.QueryRowContext(context.Background(), args...)
}

// QueryRowContext executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return sql.ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards the rest.
func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *Row {
	rows, err := s.QueryContext(ctx, args...)
	if err != nil {
		return &Row{err: err}
	}
	return &Row{rows: rows}
}
