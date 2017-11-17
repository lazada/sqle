package sqle

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/lazada/sqle/testdata"
)

var (
	origDB    *sql.DB
	db        *DB
	userId    int64
	maxUserId int64 = 1000
)

func nextUserId() (id int64) {
	if id = atomic.AddInt64(&userId, 1); id > maxUserId {
		id = 1
		atomic.StoreInt64(&userId, id)
	}
	return
}

func init() {
	var (
		w   interface{}
		err error
	)
	log.SetOutput(os.Stderr)

	if origDB, err = testdata.Open(); err != nil {
		log.Fatal(`testdata.Open() failed: `, err)
	}
	if w, err = Wrap(origDB); err != nil {
		log.Fatalf("Wrap(%T) failed: %s", origDB, err)
	}
	db = w.(*DB)

	if err = testdata.Populate(context.Background(), origDB, maxUserId); err != nil {
		log.Fatal(`testdata.Populate() failed: `, err)
	}
}
