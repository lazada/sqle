// +build go1.9

package sqle

import (
	"context"
	"testing"

	"github.com/lazada/sqle/testdata"
)

func TestWrap_Conn(t *testing.T) {
	conn, err := origDB.Conn(context.Background())
	if err != nil {
		t.Fatal(`(*sql.DB).Conn() failed:`, err)
	}
	defer conn.Close()

	w, err := Wrap(conn)
	if err != nil {
		t.Fatalf("Wrap(%T) failed: %s", conn, err)
	}
	user := testdata.User{}
	row := w.(*Conn).QueryRowContext(context.Background(), testdata.SelectUserLimitStmt, 1)
	if err = row.Scan(&user); err != nil {
		t.Fatalf("(%T).Scan() failed: %s", row, err)
	}
	if origDB != w.(*Conn).db.DB {
		t.Errorf("expected %p, got %p\n", origDB, w.(*Conn).db.DB)
	}
}
