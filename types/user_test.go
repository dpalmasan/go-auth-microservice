package types

import (
	"testing"
)

func TestValidatePasswordHash(t *testing.T) {
	user := User{
		email:        "test-user@gotest.cl",
		username:     "user1",
		passwordhash: "11123451",
		createDate:   "100000000",
		role:         0,
	}

	if !user.ValidatePasswordHash("11123451") {
		t.Fatalf(`Expected hash "11123451" to be valid`)
	}

	if user.ValidatePasswordHash("123") {
		t.Fatalf(`Expected hash "123" to be invalid`)
	}
}
