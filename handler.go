package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// authentification Middleware
func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("userName")
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

func handlerGetTracks(c *gin.Context) {
	var loadedTracks, err = GetAllTracks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tracks": loadedTracks})
}

func handlerCreateTrack(c *gin.Context) {
	var track Track
	if err := c.ShouldBindJSON(&track); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	log.Print(&track)
	id, err := Create(&track)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"_id": id})
}

func handlerGetOneTrack(c *gin.Context) {
	name := c.Param("name")

	track, err := GetOneTrack(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tracks": track})
}

// login is a handler that parses a form and checks for specific data
func login(c *gin.Context) {
	session := sessions.Default(c)
	var credentials Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	// Validate if the input is empty
	if strings.Trim(credentials.Username, " ") == "" || strings.Trim(credentials.Password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	credDB, _ := GetCredentials(&credentials)
	// check if the user is saved in the database
	if credDB == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Check for username and password match, usually from a database
	if comparePasswords(credDB["password"].(string), []byte(credentials.Password)) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Save the username and ID in the session
	session.Set("userName", credDB["username"].(string))
	session.Set("userID", credDB["_id"].(primitive.ObjectID).Hex())
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("userName")
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete("userName")
	session.Delete("_id")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func me(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("userName")
	userID := session.Get("userID")
	c.JSON(http.StatusOK, gin.H{"user": user, "ID": userID})
}

func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}

//temp funktion for generting pw hashes
func pwgen(c *gin.Context) {
	var credentials Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	pw := hashAndSalt([]byte(credentials.Password))
	c.JSON(http.StatusOK, gin.H{"Hash": pw})
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// return the hash as a string for better storage
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func handlerGetEvents(c *gin.Context) {
	var loadedEvents, err = GetAllEvents()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": loadedEvents})
}

func handlerGetOneEvent(c *gin.Context) {
	name := c.Param("name")

	event, err := GetOneEvent(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": event})
}

func handlerCreateEvent(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	log.Print(&event)
	id, err := CreateEvent(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"_id": id})
}
