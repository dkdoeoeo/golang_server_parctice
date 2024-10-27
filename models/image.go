package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID         primitive.ObjectID `bson:"_id"`
	Id         int                `bson:"id"`
	Url        string             `bson:"url"`
	Width      int                `bson:"width"`
	Height     int                `bson:"height"`
	Created_at string             `bson:"created_at"`
}

func InsertImage(image Image) error {
	collection := Mongo.Collection("image")

	result, err := collection.InsertOne(context.Background(), image)
	if err != nil {
		return err
	}

	fmt.Println("Inserted document ID:", result.InsertedID)
	return nil
}
