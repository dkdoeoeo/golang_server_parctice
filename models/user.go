package models

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func Search_user_post(c *gin.Context, user_id int, order_by, order_type string) ([]PostResponse, int, error) {
	var sortOrder bson.D
	curUser, err := GetUserById(user_id)
	if err != nil {
		log.Println("GetUserById錯誤:", err)
		return nil, 0, err
	}
	//author=user_id的post
	filter := bson.M{"author": *curUser}
	//設置排序
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

	collection := Mongo.Collection("post")
	cursor, err := collection.Find(context.Background(), filter, options.Find().SetSort(sortOrder))
	if err != nil {
		log.Println("查找post錯誤:", err)
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var posts []Post
	totalCount := 0
	//將篩選出的post解碼到posts
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

	var postResponses []PostResponse
	//將posts轉成postResponses
	for _, post := range posts {
		tmpPostResponse := new(PostResponse)
		tmpPostResponse, err = Convert_Post_To_postResponse(c, post)
		if err != nil {
			log.Println("轉換錯誤", err)
		}
		postResponses = append(postResponses, *tmpPostResponse)
	}
	return postResponses, totalCount, nil
}

func GET_user_profile(c *gin.Context, user_id int) (*User, error) {
	curUser, err := GetUserById(user_id)
	if err != nil {
		log.Println("GetUserById錯誤:", err)
		return nil, err
	}
	return curUser, err
}

func Adjust_user_profile(c *gin.Context, user_id int, nickname, profile_image string) (*User, error) {
	tmp, err := GetUserById(user_id)
	if err != nil {
		fmt.Println("查詢更新使用者失敗:", err)
		return nil, err
	}
	tokenString := c.GetHeader("Authorization")
	tokenString = tokenString[len("Bearer "):]
	if tmp.Access_token != tokenString {
		fmt.Println("非使用者本人:", err)
		return nil, err
	}

	if nickname == "" && profile_image == "" {
		fmt.Println("無須更新")
		return tmp, nil
	}

	var updatedUser User
	var updatedFields primitive.M
	collection := Mongo.Collection("user")
	filter := bson.M{"id": user_id}

	if nickname != "" && profile_image != "" {
		updatedFields = bson.M{
			"nickname":      nickname,
			"profile_image": profile_image,
		}
	} else if nickname != "" {
		updatedFields = bson.M{
			"nickname": nickname,
		}
	} else if profile_image != "" {
		updatedFields = bson.M{
			"profile_image": profile_image,
		}
	}

	update := bson.M{"$set": updatedFields}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("更新使用者失敗:", err)
		return nil, err
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	err = collection.FindOne(context.Background(), filter).Decode(&updatedUser)
	if err != nil {
		fmt.Println("查詢更新使用者失敗:", err)
		return nil, err
	}

	return &updatedUser, nil
}
