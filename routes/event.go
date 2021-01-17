package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/polaceka/Akira-API/database"
	"github.com/polaceka/Akira-API/model"
)

// GetEvents - Handler GET /event
func GetEvents(c *gin.Context) {
	var loadedEvents, err = db.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": loadedEvents})
}

// GetOneEvent - Handler GET /event/:id
func GetOneEvent(c *gin.Context) {
	id := c.Param("id")

	event, err := db.GetOneEvent(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": event})
}

// CreateEvent - Handler POST /event
func CreateEvent(c *gin.Context) {
	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	log.Print(&event)
	id, err := db.CreateEvent(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
