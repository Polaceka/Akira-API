package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionID", store))
	r.POST("/login", login)
	r.GET("/logout", logout)
	r.POST("/gen", pwgen)

	// Routing API V1
	v1 := r.Group("/v1")
	v1.Use(authRequired)
	{
		// Tacks
		v1.GET("/track", handlerGetTracks)
		v1.GET("/track/:name", handlerGetOneTrack)
		v1.POST("/track", handlerCreateTrack)

		// Events
		v1.GET("/event", handlerGetEvents)
		v1.GET("/event/:name", handlerGetOneEvent)
		v1.POST("/event", handlerCreateEvent)
		/*
			v1.Get("/event/:name/*groupe") ?
		*/
	}

	// currently for dev ussage
	r.GET("/me", authRequired, me)
	r.GET("/status", authRequired, status)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
