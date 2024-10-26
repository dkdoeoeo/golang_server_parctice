package models

type Comment struct {
	Id         int    `bson:"id"`
	User       User   `bson:"user"`
	Content    string `bson:"content"`
	Updated_at string `bson:"updated_at"`
	Created_at string `bson:"created_at"`
}
