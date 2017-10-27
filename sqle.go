package sqle

type Numerable interface {
	Num() int
}

type Aliases interface {
	Aliases() []string
	Numerable
}

type Pointers interface {
	Pointers([]interface{}, []string) ([]interface{}, int)
	Numerable
}

type Values interface {
	Values([]interface{}, []string) ([]interface{}, int)
	Numerable
}
