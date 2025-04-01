package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Online   bool               `bson:"online"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}
