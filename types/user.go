package types

type User struct {
	email        string
	username     string
	passwordhash string
	createDate   string
	role         int
}

func (u *User) ValidatePasswordHash(passHash string) bool {
	return u.passwordhash == passHash

}
