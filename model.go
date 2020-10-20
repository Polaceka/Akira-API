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
