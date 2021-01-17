package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// Track - Model of a basic track
type Track struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
}

// Credentials - User Credentials for the login
type Credentials struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Username string             `json:"username"`
	Password string             `json:"password"`
}

// Event - Model of a basic track
type Event struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	TrackName   string             `json:"trackname"`
	StartDate   primitive.DateTime `json:"stardate"`
	EndDate     primitive.DateTime `json:"enddate"`
	Group       []struct {
		Name        string `json:"name"`
		JourneyTime []struct {
			Start primitive.DateTime `json:"start"`
			End   primitive.DateTime `json:"end"`
		} `json:"journeytime"`
	} `json:"group"`
}
