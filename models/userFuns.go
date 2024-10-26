package models

import (
	"context"
	"log"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
