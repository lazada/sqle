// Generated with github.com/lazada/sqle. Do not edit by hand.
package testdata 

import (
	"github.com/lazada/sqle/embed"
)

var (
	_user_aliases_ = []string{"id", "name", "email", "created", "updated"}
	_group_aliases_ = []string{"id", "name"}
	
)

func (u *User) TypeAlias() string { return "" }

func (u *User) Aliases() []string { return _user_aliases_ }

func (u *User) Num() int { return 5 }

func (u *User) Pointers(dest []interface{}, aliases []string) (_ []interface{}, miss int) {
	if len(aliases) == 0 {
		dest = append(dest, &u.Id, &u.Name, &u.Email, &u.Created, &u.Updated)
		return dest, 0
	}
	for _, alias := range aliases {
		switch alias {
		case "id":
			dest = append(dest, &u.Id)
		case "name":
			dest = append(dest, &u.Name)
		case "email":
			dest = append(dest, &u.Email)
		case "created":
			dest = append(dest, &u.Created)
		case "updated":
			dest = append(dest, &u.Updated)
		default:
			dest, miss = append(dest, new(embed.DummyField)), miss+1
		}
	}
	return dest, miss
}

func (u *User) Values(dest []interface{}, aliases []string) ([]interface{}, error) {
	if len(aliases) == 0 {
		dest = append(dest, u.Id, u.Name, u.Email, u.Created, u.Updated)
		return dest, nil
	}
	for _, alias := range aliases {
		switch alias {
		case "id":
			dest = append(dest, u.Id)
		case "name":
			dest = append(dest, u.Name)
		case "email":
			dest = append(dest, u.Email)
		case "created":
			dest = append(dest, u.Created)
		case "updated":
			dest = append(dest, u.Updated)
		default:
			return nil, embed.ErrValueNotFound
		}
	}
	return dest, nil
}

func (g *Group) TypeAlias() string { return "" }

func (g *Group) Aliases() []string { return _group_aliases_ }

func (g *Group) Num() int { return 2 }

func (g *Group) Pointers(dest []interface{}, aliases []string) (_ []interface{}, miss int) {
	if len(aliases) == 0 {
		dest = append(dest, &g.Id, &g.Name)
		return dest, 0
	}
	for _, alias := range aliases {
		switch alias {
		case "id":
			dest = append(dest, &g.Id)
		case "name":
			dest = append(dest, &g.Name)
		default:
			dest, miss = append(dest, new(embed.DummyField)), miss+1
		}
	}
	return dest, miss
}

func (g *Group) Values(dest []interface{}, aliases []string) ([]interface{}, error) {
	if len(aliases) == 0 {
		dest = append(dest, g.Id, g.Name)
		return dest, nil
	}
	for _, alias := range aliases {
		switch alias {
		case "id":
			dest = append(dest, g.Id)
		case "name":
			dest = append(dest, g.Name)
		default:
			return nil, embed.ErrValueNotFound
		}
	}
	return dest, nil
}
