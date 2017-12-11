// `sqle` package is an extension of the standard `database/sql` package
//
// Features
//  1) fully implements the `database/sql` interface.
//  2) it is fast, sometimes very fast (has a minimum overhead).
//  3) `Scan` method can take composite types as arguments, such as structures (including nested ones), maps and slices.
//  4) `Columns` and `ColumnTypes` methods cache the returned result.
//
// Installation
//	go get -u github.com/lazada/sqle
//
// Usage
// import sql "github.com/lazada/sqle"
//
//  type User struct {
//		Id      int32     `sql:"id"`
//		Name    string    `sql:"name"`
//		Email   string    `sql:"email"`
//		Created time.Time `sql:"created"`
//		Updated time.Time `sql:"updated"`
//  }
//
//  db, err := sql.Open(`sqlite3`, `testdata/testdata.db`)
//  if err != nil {
//	  log.Fatalln(err)
//  }
//  user := new(User)
//
//  db.QueryRowContext(ctx, `SELECT * FROM users WHERE id = ?`, userId).Scan(user)
package sqle
