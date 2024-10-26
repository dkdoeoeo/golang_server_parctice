package models

type Image struct {
	Id         int    `bson:"id"`
	Url        string `bson:"url"`
	Width      int    `bson:"width"`
	Height     int    `bson:"height"`
	Created_at string `bson:"created_at"`
}
