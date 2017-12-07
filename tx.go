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

// Tx is an in-progress database transaction.
//
// A transaction must end with a call to Commit or Rollback.
//
// After a call to Commit or Rollback, all operations on the
// transaction fail with ErrTxDone.
//
// The statements prepared for a transaction by calling the transaction's
// Prepare or Stmt methods are closed by the call to Commit or Rollback.
type Tx struct {
	*sql.Tx
	db *DB
}

// Prepare creates a prepared statement for use within a transaction.
//
// The returned statement operates within the transaction and can no longer
// be used once the transaction has been committed or rolled back.
//
// To use an existing prepared statement on this transaction, see Tx.Stmt.
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	return tx.PrepareContext(context.Background(), query)
}

// Prepare creates a prepared statement for use within a transaction.
//
// The returned statement operates within the transaction and will be closed
// when the transaction has been committed or rolled back.
//
// To use an existing prepared statement on this transaction, see Tx.Stmt.
//
// The provided context will be used for the preparation of the context, not
// for the execution of the returned statement. The returned statement
// will run in the transaction context.
func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := tx.Tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{Stmt: stmt, db: tx.db}, nil
}

// Query executes a query that returns rows, typically a SELECT.
func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error) {
	return tx.QueryContext(context.Background(), query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	rows, err := tx.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: rows, db: tx.db}, nil
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
func (tx *Tx) QueryRow(query string, args ...interface{}) *Row {
	return tx.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return &Row{err: err}
	}
	return &Row{rows: rows}
}

// Stmt returns a transaction-specific prepared statement from an existing statement.
// The returned statement operates within the transaction and will be closed
// when the transaction has been committed or rolled back.
func (tx *Tx) Stmt(stmt *sql.Stmt) *Stmt {
	return tx.StmtContext(context.Background(), stmt)
}

// StmtContext returns a transaction-specific prepared statement from an existing statement.
// The returned statement operates within the transaction and will be closed
// when the transaction has been committed or rolled back.
func (tx *Tx) StmtContext(ctx context.Context, stmt *sql.Stmt) *Stmt {
	return &Stmt{Stmt: tx.Tx.StmtContext(ctx, stmt), db: tx.db}
}
