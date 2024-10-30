package models

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Id            int                `bson:"id"`
	Email         string             `bson:"email"`
	Nickname      string             `bson:"nickname"`
	Profile_image string             `bson:"profile_image"`
	Type          string             `bson:"type"`
	Access_token  string             `bson:"access_token"`
	Password      string             `bson:"password"`
	Favorite      []Post             `bson:"favorite"`
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

var (
	// 用來檢查 email 格式的正則表達式
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func IsValidImage(profile_image string) bool {
	if len(profile_image) < 4 {
		return false
	}
	ext := profile_image[len(profile_image)-4:]
	return ext == ".png" || ext == ".jpg"
}

func IsEmailUnique(email string) bool {
	//檢查email是否唯一
	collection := Mongo.Collection("user")
	var existingUser bson.M
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		return true
	} else if err != nil {
		log.Println("查詢Email出錯", err)
		return false
	} else {
		// email 已存在
		return false
	}
}

func GetUserByAccess_token(access_token string) (*User, error) {
	access_token = access_token[len("Bearer "):]
	user := new(User)
	err := Mongo.Collection("user").FindOne(context.Background(), bson.D{{"access_token", access_token}}).Decode(user)
	if err != nil {
		log.Println("Error finding user:", err)
		log.Println("access_token:", access_token)
	}
	return user, err
}

func GetUserById(user_id int) (*User, error) {
	user := new(User)
	err := Mongo.Collection("user").FindOne(context.Background(), bson.D{{"id", user_id}}).Decode(user)
	if err != nil {
		log.Println("Error finding user:", err)
		log.Println("id:", user_id)
	}
	return user, err
}
