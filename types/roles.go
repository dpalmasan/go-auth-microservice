package types

type UserRole uint8

const (
	Regular UserRole = iota
	Admin
)

func (r UserRole) String() string {
	switch r {
	case Regular:
		return "regular"
	case Admin:
		return "admin"
	}
	return "unknown"
}
