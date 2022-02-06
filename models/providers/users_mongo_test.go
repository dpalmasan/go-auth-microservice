package providers

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-auth-microservice/db/mongodb"
	"github.com/go-auth-microservice/types"
)

func setup() {
	os.Setenv("MONGO_URI", "mongodb://localhost/auth-db-test")
	os.Setenv("MONGO_DBNAME", "auth-db-test")
	mongodb.ConnectToMongo()
	mongodb.Session.Database("auth-db-test").Collection("users").Drop(context.TODO())
}

func tearDown() {
	if err := mongodb.Session.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func TestAdd(t *testing.T) {

	setup()
	user := types.User{
		Email:        "test-user@gotest.cl",
		Username:     "user1",
		Passwordhash: "11123451",
		CreatedAt:    time.Date(2020, 11, 14, 11, 30, 32, 0, time.UTC),
		Role:         0,
	}

	db := MongoDBUser{}
	insertedUser, err := db.Add(user)
	if insertedUser.Id == user.Id ||
		insertedUser.Username != user.Username ||
		insertedUser.Passwordhash != user.Passwordhash ||
		insertedUser.Role != user.Role {
		t.Fatalf(`Expected inserted user to have different id as user %+v != %+v`, insertedUser, user)
	}

	cnt, err := mongodb.Session.Database("auth-db-test").Collection("users").EstimatedDocumentCount(context.TODO())
	if err != nil || cnt != 1 {
		t.Fatalf(`Expected user count should be 1`)
	}

	tearDown()
}
