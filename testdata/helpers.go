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

package testdata

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Open() (*sql.DB, error) {
	return sql.Open(`sqlite3`, `testdata/testdata.db`)
}

func Populate(ctx context.Context, db *sql.DB, n int64) (err error) {
	var (
		stmt *sql.Stmt
		u    string
		i    int64
	)
	if _, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS users`); err != nil {
		return
	}
	if _, err = db.ExecContext(ctx, CreateUsersTableStmt); err != nil {
		return
	}
	if i, err = LastRowId(ctx, db, `users`); err != nil {
		return
	}
	i++

	if stmt, err = db.PrepareContext(ctx, `INSERT INTO users(id, name, email) VALUES (?, ?, ?)`); err != nil {
		return
	}
	for n += i; i < n; i++ {
		u = fmt.Sprintf("user%d", i)
		if _, err = stmt.ExecContext(ctx, i, u, u+`@lazada.com`); err != nil {
			return
		}
	}
	return stmt.Close()
}

func LastRowId(ctx context.Context, db *sql.DB, table string) (id int64, err error) {
	row := db.QueryRowContext(ctx, `select seq from sqlite_sequence where name = ?`, table)
	if err = row.Scan(&id); err == sql.ErrNoRows {
		err = nil
	}
	return
}

const (
	SelectUserStmt       = `SELECT * FROM users WHERE id = ?`
	SelectUserLimitStmt  = `SELECT * FROM users LIMIT ?`
	CreateUsersTableStmt = `CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(64) NOT NULL,
		email TEXT NULL,
		created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP);`
)
