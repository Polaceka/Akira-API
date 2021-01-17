package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionID", store))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	config.AllowCredentials = true
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}

	r.Use(cors.New(config))

	r.POST("/login", login)
	r.GET("/logout", logout)
	r.POST("/gen", pwgen)

	// Routing API V1
	v1 := r.Group("/v1")
	v1.Use(authRequired)
	{
		// Tacks
		v1.GET("/tracks", handlerGetTracks)
		v1.GET("/tracks/:id", handlerGetOneTrack)
		v1.POST("/tracks", handlerCreateTrack)

		// Events
		v1.GET("/events", handlerGetEvents)
		v1.GET("/events/:id", handlerGetOneEvent)
		v1.POST("/events", handlerCreateEvent)
		/*
			v1.Get("/event/:name/*groupe") ?
		*/
	}

	// currently for dev ussage
	r.GET("/me", authRequired, me)
	r.GET("/status", authRequired, status)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
