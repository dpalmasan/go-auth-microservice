package types

import "time"

type User struct {
	Id           string    `json:"id" bson:"_id,omitempty"`
	Email        string    `json:"email" bson:"email,omitempty"`
	Username     string    `json:"username" bson:"username"`
	Passwordhash string    `json:"password" bson:"password,omitempty"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at,omitempty"`
	Role         int       `json:"role" bson:"role,omitempty"`
}
