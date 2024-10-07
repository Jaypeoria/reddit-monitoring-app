package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

// Connect connects to MongoDB
func Connect() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	mongoClient = client
}

// InsertPost inserts a post into MongoDB
func InsertPost(post interface{}) {
	collection := mongoClient.Database("reddit_db").Collection("posts")
	_, err := collection.InsertOne(context.TODO(), post)
	if err != nil {
		log.Printf("Error inserting post: %v", err)
	}
}

// FindAllPosts retrieves all posts from MongoDB
func FindAllPosts() ([]bson.M, error) {
	collection := mongoClient.Database("reddit_db").Collection("posts")
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		return nil, err
	}
	defer cur.Close(context.TODO())

	var posts []bson.M
	for cur.Next(context.TODO()) {
		var post bson.M
		if err := cur.Decode(&post); err != nil {
			log.Printf("Error decoding post: %v", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
