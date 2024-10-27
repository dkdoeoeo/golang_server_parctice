package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Id            int      `bson:"id"`
	Author        User     `bson:"author"`
	Images        []Image  `bson:"images"`
	Like_count    int      `bson:"like_count"`
	Content       string   `bson:"content"`
	Type          string   `bson:"type"`
	Tags          []string `bson:"tags"`
	Location_name string   `bson:"location_name"`
	Liked         bool     `bson:"liked"`
	Updated_at    string   `bson:"updated_at"`
	Created_at    string   `bson:"created_at"`
}

func Return_public_post(order_by, order_type, content string) ([]Post, int, error) {

	var sortOrder bson.D
	filter := bson.M{"type": "public"}

	switch order_by {
	case "created_at":
		switch order_type {
		case "asc":
			sortOrder = bson.D{{"created_at", 1}}
		case "desc":
			sortOrder = bson.D{{"created_at", -1}}
		default:
			sortOrder = bson.D{{"created_at", -1}} // 默認降序
		}
	case "like_count":
		switch order_type {
		case "asc":
			sortOrder = bson.D{{"like_count", 1}}
		case "desc":
			sortOrder = bson.D{{"like_count", -1}}
		default:
			sortOrder = bson.D{{"like_count", -1}} // 默認降序
		}
	default:
		switch order_type {
		case "asc":
			sortOrder = bson.D{{"created_at", 1}}
		case "desc":
			sortOrder = bson.D{{"created_at", -1}}
		default:
			sortOrder = bson.D{{"created_at", -1}} // 默認創建時間降序
		}
	}

	if content != "" {
		filter["content"] = bson.M{"$regex": content, "$options": "i"}
	}

	collection := Mongo.Collection("post")
	cursor, err := collection.Find(context.Background(), filter, options.Find().SetSort(sortOrder))
	if err != nil {
		log.Println("查找post錯誤:", err)
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var posts []Post
	totalCount := 0

	for cursor.Next(context.Background()) {
		var post Post
		err := cursor.Decode(&post)
		if err != nil {
			log.Println("decoding post錯誤:", err)
			return nil, 0, err
		}
		posts = append(posts, post)
		totalCount++
	}

	return posts, totalCount, nil
}
