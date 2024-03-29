package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/polaceka/Akira-API/middleware"
	"github.com/polaceka/Akira-API/routes"
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

	// authentication
	r.POST("/login", routes.Login)
	r.GET("/logout", routes.Logout)

	// dev
	r.POST("/gen", routes.Pwgen)

	// Routing API V1
	v1 := r.Group("/v1")
	v1.Use(middleware.AuthRequired)
	{
		v1.POST("/track", routes.CreateTrack)
		v1.POST("/event", routes.CreateEvent)
		/*
			v1.Get("/event/:name/*group") ?
		*/
	}

	// Tacks
	v1.GET("/track", routes.GetTracks)
	v1.GET("/track/:id", routes.GetOneTrack)

	// Events
	v1.GET("/event", routes.GetEvents)
	v1.GET("/event/:id", routes.GetOneEvent)

	// currently for dev usage
	r.GET("/me", middleware.AuthRequired, routes.Me)
	r.GET("/status", middleware.AuthRequired, routes.Status)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
