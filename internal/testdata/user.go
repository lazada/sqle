package testdata

import (
	"time"
)

//go:generate sqle

//sql:"users"
type User struct {
	Id      int32      `sql:"id"`
	Name    string     `sql:"name"`
	Email   *string    `sql:"email"`
	Created time.Time  `sql:"created"`
	Updated *time.Time `sql:"updated"`
}

type Group struct {
	Id   int32  `sql:"id"`
	Name string `sql:"name"`
}
