package models

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type PostResponse struct {
	Id            int      `bson:"id"`
	Author        User     `bson:"author"`
	Images        []Image  `bson:"images"`
	Like_count    int      `bson:"like_count"`
	Content       string   `bson:"content"`
	Type          string   `bson:"type"`
	Tags          []string `bson:"tags"`
	Location_name string   `bson:"location_name"`
	Liked         *bool    `bson:"liked"`
	Updated_at    string   `bson:"updated_at"`
	Created_at    string   `bson:"created_at"`
}

func Return_public_post(c *gin.Context, order_by, order_type, content string, tags []string, location_name string) ([]PostResponse, int, error) {

	var sortOrder bson.D
	//取type=public的post
	filter := bson.M{"type": "public"}
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
	//篩選含有content的post
	if content != "" {
		filter["content"] = bson.M{"$regex": content, "$options": "i"}
	}
	//篩選含有tags的post
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$all": tags}
	}
	//篩選含有location_name的post
	if location_name != "" {
		filter["location_name"] = bson.M{"$regex": location_name, "$options": "i"}
	}

	//根據條件篩選post
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

func isUserLoggedIn(hashedEmail string) bool {
	if hashedEmail == "" {
		return false
	}

	hashedEmail = hashedEmail[len("Bearer "):]
	var user struct {
		Email string `bson:"email"`
	}
	filter := bson.M{"access_token": hashedEmail}
	UserCollection := Mongo.Collection("user")
	err := UserCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("沒找到用戶，未登入:", err)
			return false
		}
		// 其他錯誤處理
		fmt.Println("Error querying user:", err)
		return false
	}

	// 返回 true 表示用戶存在
	return true
}

func GetPostById(c *gin.Context, post_id int) (*Post, error) {
	post := new(Post)
	err := Mongo.Collection("post").FindOne(context.Background(), bson.D{{"id", post_id}}).Decode(post)
	if err != nil {
		log.Println("找不到貼文", err)
		return nil, err
	}
	return post, err
}

func Convert_Post_To_postResponse(c *gin.Context, post Post) (*PostResponse, error) {
	var author User
	UserCollection := Mongo.Collection("user")
	err := UserCollection.FindOne(context.Background(), bson.M{"_id": post.Author.ID}).Decode(&author)
	if err != nil {
		return nil, err
	}

	var images []Image
	for _, imgRef := range post.Images {
		var img Image
		ImageCollection := Mongo.Collection("image")
		err := ImageCollection.FindOne(context.Background(), bson.M{"_id": imgRef.ID}).Decode(&img)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	tmpPostResponse := PostResponse{
		Id:            post.Id,
		Author:        author, // 赋值查循到的user
		Images:        images, //赋值查循到的用image
		Like_count:    post.Like_count,
		Content:       post.Content,
		Type:          post.Type,
		Tags:          post.Tags,
		Location_name: post.Location_name,
		Updated_at:    post.Updated_at,
		Created_at:    post.Created_at,
	}

	tokenString := c.GetHeader("Authorization")
	if isUserLoggedIn(tokenString) {
		tmpPostResponse.Liked = &post.Liked
	} else {
		tmpPostResponse.Liked = nil
	}
	return &tmpPostResponse, err
}

func IsUserAuthorized(PostResponse Post, access_token string) bool {
	curUser, err := GetUserByAccess_token(access_token)
	if err != nil {
		log.Println("GetUserByAccess_token錯誤:", err)
	}

	if PostResponse.Type == "public" || PostResponse.Author.Id == curUser.Id {
		return true
	}
	return false
}

func Publish_post(c *gin.Context, Images []string, Type string, tags []string, content, location_name string) (*Post, error) {
	var images []Image
	for _, url := range Images {
		images = append(images, add_images_by_urls(url))
	}
	tokenString := c.GetHeader("Authorization")
	author, err := GetUserByAccess_token(tokenString)
	if err != nil {
		log.Println("GetUserByAccess_token錯誤:", err)
		return nil, err
	}
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02 15:04:05")
	post := Post{
		Id:            GetNextSequence("postId"),
		Author:        *author,
		Images:        images,
		Like_count:    0,
		Content:       content,
		Type:          Type,
		Tags:          tags,
		Location_name: location_name,
		Liked:         false,
		Updated_at:    timeString,
		Created_at:    timeString,
	}
	err = InsertPost(post)
	if err != nil {
		log.Println("InsertPost錯誤:", err)
	}
	return &post, err
}

func add_images_by_urls(url string) Image {
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02 15:04:05")
	fmt.Println("當前時間:", timeString)
	image := Image{
		ID:         primitive.NewObjectID(),
		Id:         GetNextSequence("imageId"),
		Url:        url,
		Width:      10,
		Height:     10,
		Created_at: timeString,
	}
	err := InsertImage(image)
	if err != nil {
		fmt.Println("插入image錯誤:", err)
	}
	return image
}

func InsertPost(post Post) error {
	collection := Mongo.Collection("post")

	result, err := collection.InsertOne(context.Background(), post)
	if err != nil {
		return err
	}

	fmt.Println("Inserted document ID:", result.InsertedID)
	return nil
}

func Adjust_post(c *gin.Context, post_id int, Type string, tags []string, content string) (*Post, error) {
	tmp, err := GetPostById(c, post_id)
	if err != nil {
		fmt.Println("查詢更新貼文失敗:", err)
		return nil, err
	}
	tokenString := c.GetHeader("Authorization")
	tokenString = tokenString[len("Bearer "):]
	if tmp.Author.Access_token != tokenString {
		fmt.Println("非貼文作者:", err)
		return nil, err
	}

	var updatedPost Post
	collection := Mongo.Collection("post")
	filter := bson.M{"id": post_id}
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02 15:04:05")
	fmt.Println("當前時間:", timeString)
	updatedFields := bson.M{
		"type":       Type,
		"tags":       tags,
		"content":    content,
		"updated_at": timeString,
	}
	update := bson.M{"$set": updatedFields}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("更新貼文失敗:", err)
		return nil, err
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	err = collection.FindOne(context.Background(), filter).Decode(&updatedPost)
	if err != nil {
		fmt.Println("查詢更新貼文失敗:", err)
		return nil, err
	}

	return &updatedPost, nil
}

func Delete_post(c *gin.Context, post_id int) error {
	tmp, err := GetPostById(c, post_id)
	if err != nil {
		fmt.Println("查詢更新貼文失敗:", err)
		return err
	}
	tokenString := c.GetHeader("Authorization")
	tokenString = tokenString[len("Bearer "):]
	if tmp.Author.Access_token != tokenString {
		fmt.Println("非貼文作者:", err)
		return err
	}

	collection := Mongo.Collection("post")
	filter := bson.M{"id": post_id}
	deleteResult, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	fmt.Printf("Matched %v documents and Deleted %v documents.\n", deleteResult.DeletedCount, deleteResult.DeletedCount)
	return err
}

func Favorite_post(c *gin.Context, post_id int) (bool, error) {
	curPost, err := GetPostById(c, post_id)
	if err != nil {
		fmt.Println("查詢貼文失敗:", err)
		return false, err
	}
	tokenString := c.GetHeader("Authorization")
	curUser, err := GetUserByAccess_token(tokenString)
	if err != nil {
		fmt.Println("查詢user失敗:", err)
		return false, err
	}

	exists := false
	for i, post := range curUser.Favorite {
		if post.Id == curPost.Id {
			curUser.Favorite = append(curUser.Favorite[:i], curUser.Favorite[i+1:]...)
			exists = true
			break
		}
	}
	if !exists {
		curUser.Favorite = append(curUser.Favorite, *curPost)
	}
	filter := bson.M{"id": curUser.Id} // 使用 curUser.ID 作為篩選條件
	update := bson.M{"$set": bson.M{"favorite": curUser.Favorite}}
	collection := Mongo.Collection("user")
	_, err = collection.UpdateOne(context.Background(), filter, update)
	return !exists, err
}

func GET_Favorite_post(c *gin.Context, order_by, order_type string) ([]Post, error) {
	tokenString := c.GetHeader("Authorization")
	if !isUserLoggedIn(tokenString) {
		fmt.Println("登入失敗:")
		return nil, nil
	}
	curUser, err := GetUserByAccess_token(tokenString)
	if err != nil {
		fmt.Println("GetUserByAccess_token錯誤:")
		return nil, err
	}
	flag := false
	if order_by == "asc" {
		flag = true
	}
	posts, err := sortUserFavorites(curUser.Favorite, flag)
	if err != nil {
		fmt.Println("sortUserFavorites錯誤:")
		return nil, err
	}
	return posts, err
}

func parsePostTime(post Post) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", post.Created_at)
}

func sortUserFavorites(userFavorites []Post, ascending bool) ([]Post, error) {
	sort.Slice(userFavorites, func(i, j int) bool {
		timeI, err := parsePostTime(userFavorites[i])
		if err != nil {
			fmt.Println("時間解析錯誤：", err)
			return false
		}
		timeJ, err := parsePostTime(userFavorites[j])
		if err != nil {
			fmt.Println("時間解析錯誤：", err)
			return false
		}
		if ascending {
			fmt.Println("升序排列")
			return timeI.Before(timeJ)
		}
		fmt.Println("降序排列")
		return timeI.After(timeJ)
	})
	return userFavorites, nil
}

func Search_user_post(c *gin.Context, user_id int, order_by, order_type string) ([]PostResponse, int, error) {
	var sortOrder bson.D
	curUser, err := GetUserById(user_id)
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
