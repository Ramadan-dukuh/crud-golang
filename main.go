package main

import (
	"crud-go/config"
	"crud-go/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	config.Connect()

	router.GET("/user", controller.GetAllUser)
	router.POST("/user", controller.SetUser)

	router.Run(":8000")
}