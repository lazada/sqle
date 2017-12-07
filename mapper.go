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
	"database/sql/driver"
	"errors"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"github.com/lazada/sqle/embed"
)

type ctorFunc func(unsafe.Pointer) unsafe.Pointer

type structMap struct {
	aliases []string
	fields  []structField
	ctors   []ctorFunc
}

func (m *structMap) alloc(ptr unsafe.Pointer) error {
	for _, ctor := range m.ctors {
		if ctor == nil {
			return ErrNilCtor
		}
		if ptr = ctor(ptr); ptr == nil {
			return ErrCtorRetNil
		}
	}
	return nil
}

type structField struct {
	offset    uintptr
	typ       reflect.Type
	ancestors []uintptr
}

type Mapper struct {
	tag   string
	conv  NamingConvention
	mu    sync.RWMutex
	types map[reflect.Type]*structMap
}

func NewMapper(tag string, conv NamingConvention) *Mapper {
	return &Mapper{
		tag:   tag,
		conv:  conv,
		types: make(map[reflect.Type]*structMap),
	}
}

func (m *Mapper) Tag() string { return m.tag }

func (m *Mapper) Aliases(src interface{}) ([]string, error) {
	typ, err := typeCheck(src)
	if err != nil {
		return nil, err
	}
	smap := m.inspect(nil, 0, typ)
	return smap.aliases, nil
}

func (m *Mapper) Pointers(src interface{}, dest []interface{}, aliases []string) ([]interface{}, int, error) {
	typ, err := typeCheck(src)
	if err != nil {
		return nil, 0, err
	}
	smap := m.inspect(nil, 0, typ)
	ptr := unsafe.Pointer(*(**byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&src)) + wordSize)))

	if len(aliases) == 0 {
		for _, f := range smap.fields {
			dest = append(dest, reflect.NewAt(f.typ, unsafe.Pointer(uintptr(ptr)+f.offset)).Interface())
		}
		return dest, 0, nil
	}
	miss := 0

	for i, j, n, m := 0, 0, len(aliases), len(smap.aliases); i < n; i++ {
		for j = 0; j < m; j++ {
			if smap.aliases[j] == aliases[i] {
				dest = append(dest, reflect.NewAt(smap.fields[j].typ, unsafe.Pointer(uintptr(ptr)+smap.fields[j].offset)).Interface())
				break
			}
		}
		if j == m {
			dest = append(dest, new(embed.DummyField))
			miss++
		}
	}
	return dest, miss, nil
}

func (m *Mapper) inspect(parent *structMap, offset uintptr, typ reflect.Type) *structMap {
	m.mu.RLock()
	smap, ok := m.types[typ]
	if m.mu.RUnlock(); ok {
		return smap
	}
	var (
		alias    string
		field    reflect.StructField
		fieldtyp reflect.Type
		isptr    uintptr
		s        *structMap
		n        int
		num      = typ.NumField()
	)
	smap = &structMap{
		aliases: make([]string, 0, num),
		fields:  make([]structField, 0, num),
		ctors:   make([]ctorFunc, 0, num),
	}
	for i, j := 0, 0; i < num; i++ {
		if field = typ.Field(i); field.PkgPath != `` && !field.Anonymous {
			continue
		}
		switch alias = field.Tag.Get(m.tag); alias {
		case `-`:
			continue
		case ``:
			if m.conv == nil {
				alias = field.Name
			} else {
				alias = m.conv.Name(field.Name)
			}
		}
		if offset&ptrMask != ptrMask {
			field.Offset += offset
		}
		if fieldtyp = field.Type; fieldtyp.Kind() == reflect.Ptr {
			fieldtyp, isptr = fieldtyp.Elem(), ptrMask
		}
		if fieldtyp.Kind() == reflect.Struct && !scannable(fieldtyp) {
			if s = m.inspect(smap, field.Offset|isptr, fieldtyp); s != nil {
				smap.aliases = append(smap.aliases, s.aliases...)
				smap.fields = append(smap.fields, s.fields...)
				smap.ctors = append(smap.ctors, s.ctors...)
				if isptr > 0 {
					n = len(smap.fields)
					for j = n - len(s.fields); j < n; j++ {
						smap.fields[j].ancestors = append(smap.fields[j].ancestors, field.Offset)
					}
					if len(s.ctors) == 0 {
						smap.ctors = append(smap.ctors, constructor(fieldtyp.Size(), field.Offset, nil))
					} else {
						n = len(smap.ctors)
						for j = n - len(s.ctors); j < n; j++ {
							smap.ctors[j] = constructor(fieldtyp.Size(), field.Offset, smap.ctors[j])
						}
					}
				}
			}
		} else {
			smap.aliases = append(smap.aliases, alias)
			smap.fields = append(smap.fields, structField{offset: field.Offset, typ: field.Type})
		}
	}
	if len(smap.fields) > 0 {
		m.mu.Lock()
		m.types[typ] = smap
		m.mu.Unlock()
		return smap
	}
	return nil
}

func typeCheck(src interface{}) (reflect.Type, error) {
	typ, ok := src.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(src)
	}
	if typ.Kind() != reflect.Ptr {
		return nil, ErrSrcNotPtr
	}
	if typ = typ.Elem(); typ.Kind() != reflect.Struct {
		return nil, ErrSrcNotPtrTo
	}
	return typ, nil
}

func constructor(size, offset uintptr, next ctorFunc) ctorFunc {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if ptr == nil {
			return nil
		}
		p := (**byte)(unsafe.Pointer(uintptr(ptr) + offset))
		if *p == nil {
			buf := make([]byte, size)
			*p = *(**byte)(unsafe.Pointer(&buf))
		}
		if next != nil {
			return next(unsafe.Pointer(*p))
		}
		return unsafe.Pointer(*p)
	}
}

func scannable(typ reflect.Type) bool {
	return typ == timeReflectType || typ.Implements(scannerReflectType)
}

const (
	ptrMask uintptr = 1 << 63
)

var (
	defaultMapper = NewMapper(`sql`, NewCachedConvention(new(SnakeConvention)))
	wordSize      = unsafe.Sizeof(uintptr(0))

	scannerReflectType = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	valuerReflectType  = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
	timeReflectType    = reflect.TypeOf(time.Time{})

	ErrSrcNotPtr   = errors.New(`sqle: source is not a pointer`)
	ErrSrcNotPtrTo = errors.New(`sqle: source is not a pointer to a structure`)
	ErrNilCtor     = errors.New(`sqle: constructor is nil`)
	ErrCtorRetNil  = errors.New(`sqle: constructor returned nil pointer`)
)
