package sqle

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"sync/atomic"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lazada/sqle/internal/testdata"
)

var (
	origDB    *sql.DB
	db        *DB
	userId    int64
	maxUserId int64 = 1000
	debug           = flag.Bool(`debug`, false, ``)
)

func nextUserId() int64 {
	id := atomic.AddInt64(&userId, 1)
	if id > maxUserId {
		id = 1
		atomic.StoreInt64(&userId, id)
	}
	return id
}

func init() {
	var (
		w   interface{}
		err error
	)
	log.SetOutput(os.Stderr)
	flag.Parse()

	if origDB, err = testdata.Open(); err != nil {
		log.Fatal(`testdata.Open() failed:`, err)
	}
	if w, err = Wrap(origDB); err != nil {
		log.Fatalf("Wrap(%T) failed: %s", origDB, err)
	}
	db = w.(*DB)

	if err = testdata.Populate(context.Background(), origDB, maxUserId); err != nil {
		log.Fatal(`testdata.Populate() failed:`, err)
	}
}
