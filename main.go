package main

import (
	"crud-go/config"
	"crud-go/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	config.Connect()

	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	router.GET("/user", controller.GetAllUser)
	router.GET("/user/:id", controller.GetUserByID)
	router.POST("/user", controller.SetUser)
	router.PUT("/user/:id", controller.UpdateUser)
	router.DELETE("/user/:id", controller.DeleteUser)

	router.Run(":8000")
}