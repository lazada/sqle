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

// DB is a database handle representing a pool of zero or more
// underlying connections. It's safe for concurrent use by multiple goroutines.
//
// The sql package creates and frees connections automatically; it
// also maintains a free pool of idle connections. If the database has
// a concept of per-connection state, such state can only be reliably
// observed within a transaction.
// Once DB.Begin is called, the returned Tx is bound to a single connection.
// Once Commit or Rollback is called on the transaction, that transaction's
// connection is returned to DB's idle connection pool.
// The pool size can be controlled with SetMaxIdleConns.
type DB struct {
	*sql.DB
	*dbOptions
}

// Begin starts a transaction. The default isolation level is dependent on the driver.
func (db *DB) Begin() (*Tx, error) { return db.BeginTx(context.Background(), nil) }

// BeginTx starts a transaction.
//
// The provided context is used until the transaction is committed or rolled back.
// If the context is canceled, the sql package will roll back the transaction.
// Tx.Commit will return an error if the context provided to BeginTx is canceled.
//
// The provided TxOptions is optional and may be nil if defaults should be used.
// If a non-default isolation level is used that the driver doesn't support,
// an error will be returned.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx, db: db}, nil
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
// The caller must call the statement's Close method when the statement is no longer needed.
func (db *DB) Prepare(query string) (*Stmt, error) {
	return db.PrepareContext(context.Background(), query)
}

// PrepareContext creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
// The caller must call the statement's Close method when the statement is no longer needed.
//
// The provided context is used for the preparation of the statement, not for the execution of the statement.
func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{Stmt: stmt, db: db}, nil
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (db *DB) Query(query string, args ...interface{}) (*Rows, error) {
	return db.QueryContext(context.Background(), query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: rows, db: db}, nil
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value.
// Errors are deferred until Row's Scan method is called.
func (db *DB) QueryRow(query string, args ...interface{}) *Row {
	return db.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value.
// Errors are deferred until Row's Scan method is called.
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return &Row{err: err}
	}
	return &Row{rows: rows}
}

// Open opens a database specified by its database driver name and a driver-specific data source name,
// usually consisting of at least a database name and connection information.
//
// Most users will open a database via a driver-specific connection helper function that returns a *DB.
// No database drivers are included in the Go standard library. See https://golang.org/s/sqldrivers for
// a list of third-party drivers.
//
// Open may just validate its arguments without creating a connection to the database.
// To verify that the data source name is valid, call Ping.
//
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once. It is rarely necessary to close a DB.
func Open(driver, dsn string, options ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	w := &DB{DB: db, dbOptions: &dbOptions{mapper: defaultMapper}}

	for _, option := range options {
		if err = option(w.dbOptions); err != nil {
			return nil, err
		}
	}
	return w, nil
}

type DBOption func(*dbOptions) error

type dbOptions struct {
	mapper *Mapper
	strict bool
}

// InStrictMode places the DB into strict mode.
func InStrictMode(opt *dbOptions) error {
	opt.strict = true
	return nil
}

// WithMapper
func WithMapper(mapper *Mapper) DBOption {
	if mapper == nil {
		mapper = defaultMapper
	}
	return func(opt *dbOptions) error {
		opt.mapper = mapper
		return nil
	}
}
