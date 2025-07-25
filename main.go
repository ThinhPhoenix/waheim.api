package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"waheim.api/configs"
	"waheim.api/handlers"
	"waheim.api/middleware"
)

func main() {
	configs.ConfDb()
	r := gin.Default()

	// CORS config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"https://thinhphoenix.github.io", "http://localhost:5173"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	auth := r.Group("/auth")
	auth.POST("/sign-up", func(c *gin.Context) {
		handlers.SignUpHandler(c.Writer, c.Request)
	})
	auth.POST("/sign-in", func(c *gin.Context) {
		handlers.SignInHandler(c.Writer, c.Request)
	})
	auth.POST("/sign-out", middleware.RequireAuthorize, func(c *gin.Context) {
		handlers.SignOutHandler(c.Writer, c.Request)
	})
	auth.GET("/me", middleware.RequireAuthorize, func(c *gin.Context) {
		handlers.AuthMeHandler(c.Writer, c.Request)
	})
	r.Run()
}
