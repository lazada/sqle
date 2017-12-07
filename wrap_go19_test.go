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
