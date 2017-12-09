# sqle
[![Build Status](https://travis-ci.org/lazada/sqle.svg?branch=master)](https://travis-ci.org/lazada/sqle)

The `sqle` package is an extension of the standard `database/sql` package.

## Features
- fully implements the `database/sql` interface.
- it is fast, sometimes very fast (has a minimum overhead).
- `Scan` method can take composite types as arguments, such as structures (including nested ones), maps and slices.
- `Columns` and `ColumnTypes` methods cache the returned result.

## Installation
```go get -u github.com/lazada/sqle```

## Usage examples
Additional examples of usage are available in [rows_test.go](https://github.com/lazada/sqle/blob/master/rows_test.go) and [row_test.go](https://github.com/lazada/sqle/blob/master/row_test.go).

### Working with structures
As it was:
```go
import "database/sql"

type User struct {
	Id      int32     `sql:"id"`
	Name    string    `sql:"name"`
	Email   string    `sql:"email"`
	Created time.Time `sql:"created"`
	Updated time.Time `sql:"updated"`
}

db, err := sql.Open(`sqlite3`, `testdata/testdata.db`)
if err != nil {
    log.Fatalln(err)
}

user := new(User)

db.QueryRowContext(ctx, `SELECT * FROM users WHERE id = ?`, userId).
    Scan(&user.Id, &user.Name, &user.Email, &user.Created, &user.Updated)
```
It is now:
```go
import sql "github.com/lazada/sqle"

// ...

user := new(User)

db.QueryRowContext(ctx, `SELECT * FROM users WHERE id = ?`, userId).Scan(user)
```
 
 ### Working with maps
As it was (simplified example):
```go
import "database/sql"

db, err := sql.Open(`sqlite3`, `testdata/testdata.db`)
if err != nil {
    log.Fatalln(err)
}
var (
	userId                   int32
	userName, userEmail      string
	userCreated, userUpdated time.Time
)
err = db.QueryRowContext(ctx, `SELECT * FROM users WHERE id = ?`, userId).
    Scan(&userId, &userName, &userEmail, &userCreated, &userUpdated)
if err != nil {
    log.Fatalln(err)
}
user := map[string]interface{}{
    `id`:      userId,
    `name`:    userName,
    `email`:   userEmail,
    `created`: userCreated,
    `updated`: userUpdated,
}
```
It is now:
```go
import sql "github.com/lazada/sqle"

db, err := sql.Open(`sqlite3`, `testdata/testdata.db`)
if err != nil {
    log.Fatalln(err)
}
user := make(map[string]interface{})

err = db.QueryRowContext(ctx, `SELECT * FROM users WHERE id = ?`, userId).Scan(user)
if err != nil {
    log.Fatalln(err)
}
```
