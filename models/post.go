package models

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
