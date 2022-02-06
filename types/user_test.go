package types

import (
	"testing"
	"time"
)

func TestValidatePasswordHash(t *testing.T) {
	user := User{
		Email:        "test-user@gotest.cl",
		Username:     "user1",
		Passwordhash: "11123451",
		CreatedAt:    time.Date(2020, 11, 14, 11, 30, 32, 0, time.UTC),
		Role:         0,
	}

	if !user.ValidatePasswordHash("11123451") {
		t.Fatalf(`Expected hash "11123451" to be valid`)
	}

	if user.ValidatePasswordHash("123") {
		t.Fatalf(`Expected hash "123" to be invalid`)
	}
}
