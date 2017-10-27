package internal

import (
	"strings"
	"text/template"
	"unicode"
)

type FileInfo struct {
	Package string
	Struct  []*StructInfo
}

type FieldInfo struct {
	Name  string
	Alias string
}

type StructInfo struct {
	Name  string
	Alias string
	Field []*FieldInfo
}

func (s *StructInfo) Aliases() []string {
	res := make([]string, len(s.Field))
	for i, f := range s.Field {
		res[i] = f.Alias
	}
	return res
}

func ShortName(name string) string {
	for _, r := range name {
		if unicode.IsLetter(r) {
			return string(unicode.ToLower(r))
		}
	}
	return `x`
}

var (
	FileTemplate = template.Must(template.New(`file`).Funcs(template.FuncMap{
		`short`: ShortName,
		`lower`: strings.ToLower,
	}).Parse(`// Generated with github.com/lazada/sqle. Do not edit by hand.
package {{.Package}} 

import (
	"github.com/lazada/sqle/embed"
)

var (
	{{range .Struct -}}
	_{{lower .Name}}_aliases_ = {{printf "%#v" .Aliases}}
	{{end}}
)
`))
	StructTemplate = template.Must(template.New(`struct`).Funcs(template.FuncMap{
		`short`: ShortName,
		`lower`: strings.ToLower,
	}).Parse(`{{$short := short .Name}}
func ({{$short}} *{{.Name}}) TypeAlias() string { return {{printf "%q" .Alias}} }

func ({{$short}} *{{.Name}}) Aliases() []string { return _{{lower .Name}}_aliases_ }

func ({{$short}} *{{.Name}}) Num() int { return {{len .Field}} }

func ({{$short}} *{{.Name}}) Pointers(dest []interface{}, aliases []string) (_ []interface{}, miss int) {
	if len(aliases) == 0 {
		dest = append(dest, {{range $i, $f := .Field}}{{if $i}}, {{end}}&{{$short}}.{{$f.Name}}{{end}})
		return dest, 0
	}
	for _, alias := range aliases {
		switch alias {
		{{range .Field -}}
		case {{printf "%q" .Alias}}:
			dest = append(dest, &{{$short}}.{{.Name}})
		{{end -}}
		default:
			dest, miss = append(dest, new(embed.DummyField)), miss+1
		}
	}
	return dest, miss
}

func ({{$short}} *{{.Name}}) Values(dest []interface{}, aliases []string) ([]interface{}, error) {
	if len(aliases) == 0 {
		dest = append(dest, {{range $i, $f := .Field}}{{if $i}}, {{end}}{{$short}}.{{$f.Name}}{{end}})
		return dest, nil
	}
	for _, alias := range aliases {
		switch alias {
		{{range .Field -}}
		case {{printf "%q" .Alias}}:
			dest = append(dest, {{$short}}.{{.Name}})
		{{end -}}
		default:
			return nil, embed.ErrValueNotFound
		}
	}
	return dest, nil
}
`))
)
