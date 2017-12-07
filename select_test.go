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
	"testing"
	"time"

	"github.com/lazada/sqle/testdata"
)

func selectStruct(b *testing.B, limit int) {
	var (
		users []*testdata.User
		u     *testdata.User
		rows  *Rows
		err   error
	)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rows, err = db.Query(testdata.SelectUserLimitStmt, limit)
		if err != nil {
			b.Fatalf("(%T).Query() failed: %s", db, err)
		}
		for rows.Next() {
			u = new(testdata.User)
			if err = rows.Scan(u); err != nil {
				b.Errorf("(%T).Scan() failed: %s", rows, err)
			}
			users = append(users, u)
		}
	}
}
func BenchmarkSelect1_Struct(b *testing.B)    { selectStruct(b, 1) }
func BenchmarkSelect10_Struct(b *testing.B)   { selectStruct(b, 10) }
func BenchmarkSelect100_Struct(b *testing.B)  { selectStruct(b, 100) }
func BenchmarkSelect1000_Struct(b *testing.B) { selectStruct(b, 1000) }

func selectAnonStruct(b *testing.B, limit int) {
	var (
		users []*struct {
			Id      int
			Name    string
			Email   *string
			Created time.Time
			Updated *time.Time
		}
		u *struct {
			Id      int
			Name    string
			Email   *string
			Created time.Time
			Updated *time.Time
		}
		rows *Rows
		err  error
	)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rows, err = db.Query(testdata.SelectUserLimitStmt, limit)
		if err != nil {
			b.Fatalf("(%T).Query() failed: %s", db, err)
		}
		for rows.Next() {
			u = new(struct {
				Id      int
				Name    string
				Email   *string
				Created time.Time
				Updated *time.Time
			})
			if err = rows.Scan(u); err != nil {
				b.Errorf("(%T).Scan() failed: %s", rows, err)
			}
			users = append(users, u)
		}
	}
}
func BenchmarkSelect1_AnonStruct(b *testing.B)    { selectAnonStruct(b, 1) }
func BenchmarkSelect10_AnonStruct(b *testing.B)   { selectAnonStruct(b, 10) }
func BenchmarkSelect100_AnonStruct(b *testing.B)  { selectAnonStruct(b, 100) }
func BenchmarkSelect1000_AnonStruct(b *testing.B) { selectAnonStruct(b, 1000) }

func selectMap(b *testing.B, limit int) {
	var (
		users []map[string]interface{}
		u     map[string]interface{}
		rows  *Rows
		err   error
	)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rows, err = db.Query(testdata.SelectUserLimitStmt, limit)
		if err != nil {
			b.Fatalf("(%t).Query() failed: %s", db, err)
		}
		for rows.Next() {
			u = make(map[string]interface{})
			if err = rows.Scan(u); err != nil {
				b.Errorf("(%T).Scan() failed:", rows, err)
			}
			users = append(users, u)
		}
	}
}
func BenchmarkSelect1_Map(b *testing.B)    { selectMap(b, 1) }
func BenchmarkSelect10_Map(b *testing.B)   { selectMap(b, 10) }
func BenchmarkSelect100_Map(b *testing.B)  { selectMap(b, 100) }
func BenchmarkSelect1000_Map(b *testing.B) { selectMap(b, 1000) }
