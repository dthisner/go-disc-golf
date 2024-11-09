package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Snippet struct {
	Snippet_id int       `bson:"snippet_id"`
	Title      string    `bson:"title"`
	Content    string    `bson:"content"`
	Created    time.Time `bson:"created"`
	Expires    time.Time `bson:"expires"`
}

type SnippetModel struct {
	DB *mongo.Client
}

// Insert This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	collection := m.DB.Database("snippetsDB").Collection("snippets")

	var latestSnippet Snippet
	err := collection.FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.D{{Key: "created", Value: -1}})).Decode(&latestSnippet)
	if err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	}
	newID := latestSnippet.Snippet_id + 1
	log.Println("NewId", newID)

	snippet := Snippet{
		Snippet_id: newID,
		Title:      title,
		Content:    content,
		Created:    time.Now(),
		Expires:    time.Now().Add(time.Duration(expires) * time.Second),
	}

	_, err = collection.InsertOne(context.Background(), snippet)
	if err != nil {
		log.Printf("Failed to insert snippet: %v", err)
		return 0, err
	}

	return newID, nil
}

// Get This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	collection := m.DB.Database("snippetsDB").Collection("snippets")
	var snippet Snippet

	filter := bson.D{{Key: "snippet_id", Value: id}}
	err := collection.FindOne(context.Background(), filter).Decode(&snippet)
	if err == mongo.ErrNoDocuments {
		return Snippet{}, ErrNoRecord
	} else if err != nil {
		return Snippet{}, err
	}

	return snippet, nil
}

// Latest This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	collection := m.DB.Database("snippetsDB").Collection("snippets")
	var snippets []Snippet

	// Create a sort filter (sort by 'createdAt' field in descending order)
	sortOptions := bson.D{{Key: "created", Value: -1}} // Descending order (-1)

	findOptions := options.Find()
	findOptions.SetSort(sortOptions)
	findOptions.SetLimit(10)

	cursor, err := collection.Find(context.Background(), bson.D{}, findOptions)
	if err == mongo.ErrNoDocuments {
		return []Snippet{}, ErrNoRecord
	} else if err != nil {
		return []Snippet{}, err
	}

	// Iterate through the cursor and decode each document into the Snippet struct
	for cursor.Next(context.Background()) {
		var snippet Snippet
		if err := cursor.Decode(&snippet); err != nil {
			log.Fatal(err)
		}
		snippets = append(snippets, snippet)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return snippets, nil
}
