package strcase

import "testing"

func TestSnake(t *testing.T) {
	cases := [...]struct {
		in  string
		exp []string
	}{
		{` abc  abc   `, []string{`abc_abc`, `_abc_abc_`, `_abc__abc__`, `_abc__abc___`}},
		{`_abc__abc___`, []string{`abc_abc`, `_abc_abc_`, `_abc__abc__`, `_abc__abc___`}},
		{`-abc--abc---`, []string{`abc_abc`, `_abc_abc_`, `_abc__abc__`, `_abc__abc___`}},
		{` абв  абв   `, []string{`абв_абв`, `_абв_абв_`, `_абв__абв__`, `_абв__абв___`}},
		{` абв- абв_  `, []string{`абв_абв`, `_абв_абв_`, `_абв__абв__`, `_абв__абв___`}},
	}
	for _, c := range cases {
		for g := 0; g < len(c.exp); g++ {
			if out := Snake(c.in, '_', uint8(g)); c.exp[g] != out {
				t.Errorf("for: %q, exp: %q, got: %q", c.in, c.exp[g], out)
			}
		}
	}
}

func TestToSnake(t *testing.T) {
	cases := [...]struct{ in, exp string }{
		{``, ``},
		{`abcabc`, `abcabc`},
		{`ABCABC`, `abcabc`},
		{`Abcabc`, `abcabc`},
		{`abcAbc`, `abc_abc`},
		{`AbcAbc`, `abc_abc`},
		{`ABCAbc`, `abc_abc`},
		{`abcABC`, `abc_abc`},
		{`abc_abc`, `abc_abc`},
		{`abc_Abc`, `abc_abc`},
		{`AbCaBc`, `ab_ca_bc`},
		{`абвабв`, `абвабв`},
		{`АБВАБВ`, `абвабв`},
		{`Абвабв`, `абвабв`},
		{`абвАбв`, `абв_абв`},
		{`АбвАбв`, `абв_абв`},
		{`АБВАбв`, `абв_абв`},
		{`абвАБВ`, `абв_абв`},
	}
	for _, c := range cases {
		if out := ToSnake(c.in); c.exp != out {
			t.Errorf("for: %q, exp: %q, got: %q", c.in, c.exp, out)
		}
	}
}

func TestToKebab(t *testing.T) {
	cases := [...]struct{ in, exp string }{
		{``, ``},
		{`abcabc`, `abcabc`},
		{`ABCABC`, `abcabc`},
		{`Abcabc`, `abcabc`},
		{`abcAbc`, `abc-abc`},
		{`AbcAbc`, `abc-abc`},
		{`ABCAbc`, `abc-abc`},
		{`abcABC`, `abc-abc`},
		{`abc-abc`, `abc-abc`},
		{`AbCaBc`, `ab-ca-bc`},
	}
	for _, c := range cases {
		if out := ToKebab(c.in); c.exp != out {
			t.Errorf("for: %q, exp: %q, got: %q", c.in, c.exp, out)
		}
	}
}
