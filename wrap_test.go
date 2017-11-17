package sqle

import (
	"context"
	"testing"

	"github.com/lazada/sqle/testdata"
)

func TestWrap_DB(t *testing.T) {
	w, err := Wrap(origDB)
	if err != nil {
		t.Fatalf("Wrap(%T) failed: %s", origDB, err)
	}
	if err = w.(*DB).Ping(); err != nil {
		t.Fatal(`Ping() failed:`, err)
	}
	if origDB != w.(*DB).DB {
		t.Errorf("expected %p, got %p", origDB, w.(*DB).DB)
	}
}

func TestWrap_Rows(t *testing.T) {
	rows, err := origDB.Query(testdata.SelectUserLimitStmt, 1)
	if err != nil {
		t.Fatalf("(%T).Query() failed: %s", origDB, err)
	}
	defer rows.Close()

	w, err := Wrap(rows)
	if err != nil {
		t.Fatalf("Wrap(%T) failed: %s", rows, err)
	}
	user := testdata.User{}
	w.(*Rows).Next()

	if err = w.(*Rows).Scan(&user); err != nil {
		t.Fatalf("(%T).Scan() failed: %s", w.(*Rows), err)
	}
	if origDB != w.(*Rows).db.DB {
		t.Errorf("expected %p, got %p", origDB, w.(*Rows).db.DB)
	}
}

func TestWrap_ClosedRows(t *testing.T) {
	rows, err := origDB.Query(testdata.SelectUserLimitStmt, 1)
	if err != nil {
		t.Fatalf("(%T).Query() failed: %s", origDB, err)
	}
	if err = rows.Close(); err != nil {
		t.Fatalf("(%T).Close() failed: %s", rows, err)
	}
	w, err := Wrap(rows)
	if err != nil {
		t.Fatalf("Wrap(%T) failed: %s", rows, err)
	}
	if origDB != w.(*Rows).db.DB {
		t.Errorf("expected %p, got %p", origDB, w.(*Rows).db.DB)
	}
}

func TestWrap_Row(t *testing.T) {
	w, err := Wrap(origDB.QueryRow(testdata.SelectUserLimitStmt, 1))
	if err != nil {
		t.Fatal(`Wrap(*sql.Row) failed:`, err)
	}
	user := testdata.User{}
	if err = w.(*Row).Scan(&user); err != nil {
		t.Fatalf("(%T).Scan() failed: %s", w.(*Row), err)
	}
	if origDB != w.(*Row).rows.db.DB {
		t.Errorf("expected %p, got %p\n", origDB, w.(*Row).rows.db.DB)
	}
}

func TestWrap_Tx(t *testing.T) {
	tx, err := origDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("(%T).BeginTx() failed: %s", origDB, err)
	}
	w, err := Wrap(tx)
	if err != nil {
		t.Fatalf("Wrap(%T) failed: %s", tx, err)
	}
	row, user := w.(*Tx).QueryRow(testdata.SelectUserLimitStmt, 1), testdata.User{}
	if err = row.Scan(&user); err != nil {
		t.Fatalf("(%T).Scan() failed: %s", row, err)
	}
	if err = w.(*Tx).Commit(); err != nil {
		t.Fatalf("(%T).Commit() failed: %s", w.(*Tx), err)
	}
	if origDB != w.(*Tx).db.DB {
		t.Errorf("expected %p, got %p\n", origDB, w.(*Tx).db.DB)
	}
}

func TestWrap_Stmt(t *testing.T) {
	stmt, err := origDB.Prepare(testdata.SelectUserLimitStmt)
	if err != nil {
		t.Fatalf("(%T).Prepare() failed: %s", origDB, err)
	}
	defer stmt.Close()

	w, err := Wrap(stmt)
	if err != nil {
		t.Fatalf("Wrap(%T) failed: %s", stmt, err)
	}
	row, user := w.(*Stmt).QueryRow(1), testdata.User{}
	if err = row.Scan(&user); err != nil {
		t.Fatalf("(%T).Scan() failed: %s", row, err)
	}
	if origDB != w.(*Stmt).db.DB {
		t.Errorf("expected %p, got %p\n", origDB, w.(*Stmt).db.DB)
	}
}
