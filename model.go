package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// Track - Model of a basic track
type Track struct {
	_ID         primitive.ObjectID
	Name        string
	Description string
}

// Credentials - User Credentials for the login
type Credentials struct {
	_ID      primitive.ObjectID
	Username string
	Password string
}

// Event - Model of a basic track
type Event struct {
	_ID         primitive.ObjectID
	Name        string
	Description string
	TrackName   string
	StartDate   primitive.DateTime
	EndDate     primitive.DateTime
	Group       []struct {
		Name        string
		JourneyTime []struct {
			Start primitive.DateTime
			End   primitive.DateTime
		}
	}
}
