package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	db "github.com/polaceka/Akira-API/database"
	"github.com/polaceka/Akira-API/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// GetTracks -
func GetTracks(c *gin.Context) {
	var loadedTracks, err = db.GetAllTracks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tracks": loadedTracks})
}

// CreateTrack -
func CreateTrack(c *gin.Context) {
	var track model.Track
	if err := c.ShouldBindJSON(&track); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	log.Print(&track)
	id, err := db.Create(&track)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"_id": id})
}

// GetOneTrack -
func GetOneTrack(c *gin.Context) {
	id := c.Param("id")

	track, err := db.GetOneTrack(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tracks": track})
}

// Login - in is a handler that parses a form and checks for specific data
func Login(c *gin.Context) {
	session := sessions.Default(c)
	var credentials model.Credentials
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

	credDB, _ := db.GetCredentials(&credentials)
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

func Logout(c *gin.Context) {
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

func Me(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("userName")
	userID := session.Get("userID")
	c.JSON(http.StatusOK, gin.H{"user": user, "id": userID})
}

// Status -
func Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}

//temp funktion for generting pw hashes
func Pwgen(c *gin.Context) {
	var credentials model.Credentials
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
