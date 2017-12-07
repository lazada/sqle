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

package testdata

import (
	"time"
)

//go:generate sqle

//sql:"users"
type User struct {
	Id      int32      `sql:"id"`
	Name    string     `sql:"name"`
	Email   *string    `sql:"email"`
	Created time.Time  `sql:"created"`
	Updated *time.Time `sql:"updated"`
}

type Group struct {
	Id   int32  `sql:"id"`
	Name string `sql:"name"`
}
