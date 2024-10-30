package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Comment struct {
	Id         int    `bson:"id"`
	User       User   `bson:"user"`
	Content    string `bson:"content"`
	Updated_at string `bson:"updated_at"`
	Created_at string `bson:"created_at"`
	Post_id    int    `bson:"post_id"`
}

type CommentRequest struct {
	Content string `json:"content"`
}

func Post_comment(c *gin.Context, post_id int, content string) (*Comment, error) {
	tokenString := c.GetHeader("Authorization")
	curUser, err := GetUserByAccess_token(tokenString)
	if err != nil {
		fmt.Println("GetUserByAccess_token錯誤：", err)
		return nil, err
	}
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02 15:04:05")
	newComment := Comment{
		Id:         GetNextSequence("commentId"),
		User:       *curUser,
		Content:    content,
		Updated_at: timeString,
		Created_at: timeString,
		Post_id:    post_id,
	}
	collection := Mongo.Collection("comment")
	result, err := collection.InsertOne(context.Background(), newComment)
	if err != nil {
		return nil, err
	}

	fmt.Println("Inserted document ID:", result.InsertedID)
	return &newComment, nil
}

func Adjust_comment(c *gin.Context, post_id int, comment_id int, content string) (*Comment, error) {
	tmp, err := GetCommentById(comment_id)
	if err != nil {
		fmt.Println("查詢更新貼文失敗:", err)
		return nil, err
	}
	tokenString := c.GetHeader("Authorization")
	curUser, err := GetUserByAccess_token(tokenString)
	if tmp.User.Id != curUser.Id {
		fmt.Println("非留言作者，更新失敗:", err)
		return nil, err
	}

	collection := Mongo.Collection("comment")
	filter := bson.M{"id": comment_id}
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02 15:04:05")
	fmt.Println("當前時間:", timeString)
	updatedFields := bson.M{
		"content":    content,
		"updated_at": timeString,
	}
	update := bson.M{"$set": updatedFields}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("更新留言失敗:", err)
		return nil, err
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	var updatedComment Comment
	err = collection.FindOne(context.Background(), filter).Decode(&updatedComment)
	if err != nil {
		fmt.Println("查詢更新留言失敗:", err)
		return nil, err
	}

	return &updatedComment, nil
}

func GetCommentById(comment_id int) (*Comment, error) {
	comment := new(Comment)
	err := Mongo.Collection("comment").FindOne(context.Background(), bson.D{{"id", comment_id}}).Decode(comment)
	if err != nil {
		log.Println("找不到留言", err)
		return nil, err
	}
	return comment, err
}

func Delete_comment(c *gin.Context, post_id int, comment_id int) error {
	tmp, err := GetCommentById(comment_id)
	if err != nil {
		fmt.Println("查詢刪除貼文失敗:", err)
		return err
	}
	tokenString := c.GetHeader("Authorization")
	curUser, err := GetUserByAccess_token(tokenString)
	if tmp.User.Id != curUser.Id {
		fmt.Println("非留言作者，刪除失敗:", err)
		return err
	}

	collection := Mongo.Collection("comment")
	filter := bson.M{"id": comment_id}
	deleteResult, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	fmt.Printf("Matched %v documents and Deleted %v documents.\n", deleteResult.DeletedCount, deleteResult.DeletedCount)
	return err
}
