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

	// Auth routes
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
	router.POST("/user", controller.Register) // For backward compatibility with frontend

	// User routes (Admin only can manage users)
	router.GET("/user", controller.AuthMiddleware("admin"), controller.GetAllUser)
	router.GET("/user/:id", controller.GetUserByID)
	router.PUT("/user/:id", controller.AuthMiddleware("admin"), controller.UpdateUser)
	router.DELETE("/user/:id", controller.AuthMiddleware("admin"), controller.DeleteUser)

	// Barang routes
	router.GET("/barang", controller.GetAllBarang)
	router.GET("/barang/:id", controller.GetBarangByID)
	router.POST("/barang", controller.AuthMiddleware("admin", "seller"), controller.CreateBarang)
	router.PUT("/barang/:id", controller.AuthMiddleware("admin"), controller.UpdateBarang)
	router.DELETE("/barang/:id", controller.AuthMiddleware("admin", "seller"), controller.DeleteBarang)

	// Cart routes (Customer only)
	router.GET("/cart/:user_id", controller.GetCartByUserID)
	router.POST("/cart", controller.AuthMiddleware("customer"), controller.AddToCart)
	router.PUT("/cart/:cart_id", controller.UpdateCart)
	router.DELETE("/cart/:cart_id", controller.RemoveFromCart)
	router.DELETE("/cart/clear/:user_id", controller.ClearCart)

	// Orders routes
	router.GET("/orders", controller.GetAllOrders)
	router.GET("/orders/:id", controller.GetOrderByID)
	router.GET("/orders/user/:user_id", controller.GetOrdersByUserID)
	router.POST("/orders", controller.CreateOrder)
	router.PUT("/orders/:id", controller.UpdateOrderStatus)
	router.DELETE("/orders/:id", controller.DeleteOrder)

	router.Run(":8000")
}