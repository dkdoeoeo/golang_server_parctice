package test

import (
	"context"
	"fmt"
	"post-platform/models"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestFindOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://140.118.178.117:27017"))
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database("post-paltform")
	ub := new(models.User)
	err = db.Collection("user").FindOne(context.Background(), bson.D{{"email", "String"}, {"password", "string"}}).Decode(ub)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ub ==> ", ub)
}
