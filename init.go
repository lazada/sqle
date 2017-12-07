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
	"database/sql"
	"reflect"
)

func init() {
	var (
		rowtyp = reflect.TypeOf(sql.Row{})
		found  bool
	)
	for i, n := 0, rowtyp.NumField(); i < n; i++ {
		if field := rowtyp.Field(i); field.Name == `rows` {
			offsetRowRows, found = field.Offset, true
		}
	}
	if !found {
		panic(`unexpected structure of database/sql/Row`)
	}
}

var (
	offsetRowRows uintptr
)
