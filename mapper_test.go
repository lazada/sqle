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
	"testing"
	"time"
)

type testStruct struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	Price     *sql.NullFloat64
	Price2    sql.NullFloat64
}

func TestInspect(t *testing.T) {
	mapper := NewMapper("sql", nil)
	typ := reflect.TypeOf(&testStruct{}).Elem()

	smap := mapper.inspect(nil, 0, typ)

	expectedAliases := []string{"ID", "Name", "CreatedAt", "UpdatedAt", "Price", "Price2"}
	if !reflect.DeepEqual(smap.aliases, expectedAliases) {
		t.Errorf("Expected aliases %v, but got %v", expectedAliases, smap.aliases)
	}

	expectedFieldsCount := len(expectedAliases)
	if len(smap.fields) != expectedFieldsCount {
		t.Errorf("Expected %d fields, but got %d", expectedFieldsCount, len(smap.fields))
	}

	for i, field := range smap.fields {
		expectedOffset := typ.Field(i).Offset
		if field.offset != expectedOffset {
			t.Errorf("Expected offset %d for field %s, but got %d", expectedOffset, expectedAliases[i], field.offset)
		}
	}

	expectedTypes := []reflect.Type{
		reflect.TypeOf(int(0)),
		reflect.TypeOf(""),
		reflect.TypeOf(time.Time{}),
		reflect.TypeOf(&time.Time{}),
		reflect.TypeOf(&sql.NullFloat64{}),
		reflect.TypeOf(sql.NullFloat64{}),
	}
	for i, field := range smap.fields {
		if field.typ != expectedTypes[i] {
			t.Errorf("Expected type %v for field %s, but got %v", expectedTypes[i], expectedAliases[i], field.typ)
		}
	}
}
