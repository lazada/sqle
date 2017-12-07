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
