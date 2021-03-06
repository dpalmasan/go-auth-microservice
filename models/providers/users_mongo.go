package providers

import (
	"context"
	"errors"
	"time"

	"github.com/go-auth-microservice/db/mongodb"
	"github.com/go-auth-microservice/types"
	"github.com/go-auth-microservice/utils/crypto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CollectionUser = "users"
)

type MongoDBUser struct{}

func (m MongoDBUser) GetByEmail(email string) (types.User, error) {
	var user types.User
	var result bson.M
	coll := mongodb.Session.Database(mongodb.DatabaseName).Collection(CollectionUser)
	err := coll.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return user, err
	}

	if err != nil {
		panic(err)
	}
	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &user)
	return user, nil
}

func (m MongoDBUser) Add(user types.User) (types.User, error) {
	coll := mongodb.Session.Database(mongodb.DatabaseName).Collection(CollectionUser)
	var result types.User
	err := coll.FindOne(context.TODO(), bson.M{
		"$or": []bson.M{
			{"email": user.Email},
			{"username": user.Username},
		}}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		time := time.Now()
		user.CreatedAt = time
		hashedPassword, err := crypto.HashPassword(user.Passwordhash)

		if err != nil {
			return user, err
		}

		user.Passwordhash = hashedPassword

		// We always want users registered to be non admin
		user.Role = types.Regular

		res, err := mongodb.Session.Database(mongodb.DatabaseName).Collection(CollectionUser).InsertOne(nil, user)
		if err != nil {
			return user, err
		}

		user.Id = res.InsertedID.(primitive.ObjectID).Hex()

		return user, nil

	}
	if err == nil {
		if result.Username == user.Username {
			return user, errors.New("Username already exist")
		}
		return user, errors.New("Email is already registered")
	}
	return user, err
}

func (m MongoDBUser) GetById(id string) (types.User, error) {
	var user types.User
	var result bson.M

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	coll := mongodb.Session.Database(mongodb.DatabaseName).Collection(CollectionUser)
	err = coll.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return user, err
	}

	if err != nil {
		panic(err)
	}
	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &user)
	return user, nil
}
