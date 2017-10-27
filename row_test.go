package sqle

import (
	"reflect"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lazada/sqle/internal/testdata"
)

func TestRow_ScanMap(t *testing.T) {
	var (
		id               int64
		name, email      string
		created, updated string
		exp              = make(map[string]interface{})
		out              = make(map[string]interface{})
		uid              = nextUserId()
	)
	err := origDB.QueryRow(testdata.SelectUserStmt, uid).Scan(
		&id, &name, &email, &created, &updated,
	)
	if err != nil {
		t.Fatalf("(%T).QueryRow().Scan() failed: %s", origDB, err)
	}
	exp[`id`] = id
	exp[`name`], exp[`email`] = name, email
	exp[`created`], exp[`updated`] = created, updated

	if err = db.QueryRow(testdata.SelectUserStmt, uid).Scan(out); err != nil {
		t.Fatalf("(%T).QueryRow().Scan() failed: %s", db, err)
	}
	if !reflect.DeepEqual(exp, out) {
		t.Errorf("expected %v, got %v\n", exp, out)
	}
}

func TestRow_ScanPtrMap(t *testing.T) {
	var (
		id               int64
		name, email      string
		created, updated string
		exp              = make(map[string]interface{})
		out              = make(map[string]interface{})
		uid              = nextUserId()
	)
	err := origDB.QueryRow(testdata.SelectUserStmt, uid).Scan(
		&id, &name, &email, &created, &updated,
	)
	if err != nil {
		t.Fatalf("(%T).QueryRow().Scan() failed: %s", origDB, err)
	}
	exp[`id`] = id
	exp[`name`], exp[`email`] = name, email
	exp[`created`], exp[`updated`] = created, updated

	if err = db.QueryRow(testdata.SelectUserStmt, uid).Scan(&out); err != nil {
		t.Fatalf("(%T).QueryRow().Scan() failed: %s", db, err)
	}
	if !reflect.DeepEqual(exp, out) {
		t.Errorf("expected %v, got %v\n", exp, out)
	}
}

func TestRow_ScanAnonStruct(t *testing.T) {
	exp := struct {
		Id      int
		Name    string
		Email   *string
		Created time.Time
		Updated *time.Time
	}{}
	out, uid := exp, nextUserId()
	err := origDB.QueryRow(testdata.SelectUserStmt, uid).Scan(
		&exp.Id, &exp.Name, &exp.Email, &exp.Created, &exp.Updated,
	)
	if err != nil {
		t.Fatalf("(%T).QueryRow().Scan() failed: %s", origDB, err)
	}
	if err = db.QueryRow(testdata.SelectUserStmt, uid).Scan(&out); err != nil {
		t.Fatalf("(%T).QueryRow().Scan() failed: %s", db, err)
	}
	if !reflect.DeepEqual(exp, out) {
		t.Errorf("expected %v, got %v\n", exp, out)
	}
}

func TestRow_ScanVarAnonStructVar(t *testing.T) {
	id, updated := 0, new(time.Time)
	u := struct {
		Name    string
		Email   *string
		Created time.Time
	}{}
	if err := db.QueryRow(testdata.SelectUserStmt, nextUserId()).Scan(
		&id, &u, &updated,
	); err != nil {
		t.Fatalf("(%T).QueryRow().Scan() failed: %s", db, err)
	}
	// fmt.Printf("id: %d, updated: %s, struct: %#v\n", id, updated, u)
}
