package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Id            int    `bson:"id"`
	Email         string `bson:"email"`
	Nickname      string `bson:"nickname"`
	Profile_image string `bson:"profile_image"`
	Type          string `bson:"type"`
	Password      string `bson:"password"`
}

func GetUserByEmailPassword(email, password string) (*User, error) {
	user := new(User)
	err := Mongo.Collection("user").FindOne(context.Background(), bson.D{{"email", email}, {"password", password}}).Decode(user)
	if err != nil {
		log.Println("Error finding user:", err)
	}
	return user, err
}

func GetNextSequence(seqName string) int {
	collection := Mongo.Collection("counters")

	filter := bson.M{"_id": seqName}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	var result struct {
		Seq int `bson:"seq"`
	}

	err := collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&result)
	if err != nil {
		return 0
	}

	return result.Seq
}

func InsertUser(user User) error {
	collection := Mongo.Collection("user")

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	fmt.Println("Inserted document ID:", result.InsertedID)
	return nil
}
