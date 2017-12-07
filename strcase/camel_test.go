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

package strcase

import "testing"

func TestToCamel(t *testing.T) {
	cases := [...]struct{ in, exp string }{
		{``, ``},
		{`abcabc`, `Abcabc`},
		{`abcAbc`, `AbcAbc`},
		{`ABCabc`, `ABCabc`},
		{`abc abc`, `AbcAbc`},
		{`abc_abc`, `AbcAbc`},
		{`abc-abc`, `AbcAbc`},
		{`  abc  abc    `, `AbcAbc`},
		{`__abc__abc____`, `AbcAbc`},
		{`--abc--abc----`, `AbcAbc`},
		{`абвабв`, `Абвабв`},
		{`абвАбв`, `АбвАбв`},
	}
	for _, c := range cases {
		if out := ToCamel(c.in); out != c.exp {
			t.Errorf("for %q, exp: %q, got: %q", c.in, c.exp, out)
		}
	}
}
