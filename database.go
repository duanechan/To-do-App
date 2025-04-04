package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var uri string = "mongodb+srv://%s:%s@to-do.qj1dwji.mongodb.net/?retryWrites=true&w=majority&appName=To-do"

var mongoClient *mongo.Client

func init() {
	if err := connect(); err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}
}

func connect() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(fmt.Sprintf(uri, username, password)).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	err, mongoClient = client.Ping(ctx, nil), client

	return err
}

func Login(username, password string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	collection := mongoClient.Database("ToDo").Collection("Users")

	var user User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return User{}, err
	}

	user.Online = true

	_, err = collection.ReplaceOne(ctx, bson.M{"_id": user.ID, "username": user.Username}, user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
