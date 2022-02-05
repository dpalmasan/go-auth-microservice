package types

type user struct {
	email        string
	username     string
	passwordhash string
	createDate   string
	role         int
}

func (u *user) ValidatePasswordHash(passHash string) bool {
	return u.passwordhash == passHash

}
