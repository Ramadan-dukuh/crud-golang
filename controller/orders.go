package controller

import (
	"crud-go/config"
	"crud-go/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAllOrders - Get all orders with pagination
func GetAllOrders(c *gin.Context) {
	var orders []model.Orders

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "5")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 5
	}

	offset := (pageInt - 1) * limitInt

	var total int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": err.Error(),
			"status": 0,
		})
		return
	}

	rows, err := config.DB.Query("SELECT id, user_id, total_harga, status, createdAt FROM orders LIMIT ? OFFSET ?", limitInt, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": err.Error(),
			"status": 0,
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var o model.Orders
		err := rows.Scan(
			&o.Id,
			&o.User_id,
			&o.Total_harga,
			&o.Status,
			&o.Created_at,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": err.Error(),
				"status": 0,
			})
			return
		}
		orders = append(orders, o)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": orders,
		"status": 1,
		"page":   pageInt,
		"limit":  limitInt,
		"total":  total,
	})
}

// GetOrdersByUserID - Get orders for a specific user
func GetOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	var orders []model.Orders

	query := "SELECT id, user_id, total_harga, status, createdAt FROM orders WHERE user_id = ?"
	rows, err := config.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": err.Error(),
			"status": 0,
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var o model.Orders
		err := rows.Scan(
			&o.Id,
			&o.User_id,
			&o.Total_harga,
			&o.Status,
			&o.Created_at,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": err.Error(),
				"status": 0,
			})
			return
		}
		orders = append(orders, o)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": orders,
		"status": 1,
	})
}

// GetOrderByID - Get order by ID
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order model.Orders

	query := "SELECT id, user_id, total_harga, status, createdAt FROM orders WHERE id = ?"
	err := config.DB.QueryRow(query, id).Scan(
		&order.Id,
		&order.User_id,
		&order.Total_harga,
		&order.Status,
		&order.Created_at,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "Order tidak ditemukan",
			"status": 0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": order,
		"status": 1,
	})
}

// CreateOrder - Create new order
func CreateOrder(c *gin.Context) {
	var newOrder model.Orders

	err := c.BindJSON(&newOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	newOrder.Created_at = time.Now()
	if newOrder.Status == "" {
		newOrder.Status = "pending"
	}

	query := "INSERT INTO orders (user_id, total_harga, status, createdAt) VALUES (?, ?, ?, ?)"
	result, err := config.DB.Exec(query, newOrder.User_id, newOrder.Total_harga, newOrder.Status, newOrder.Created_at)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	if affected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": "Order berhasil dibuat",
			"status": 1,
			"data":   newOrder,
		})
	}
}

// UpdateOrderStatus - Update order status
func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var updateOrder model.Orders

	err := c.BindJSON(&updateOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	query := "UPDATE orders SET status = ? WHERE id = ?"
	result, err := config.DB.Exec(query, updateOrder.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	if affected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": "Status order berhasil diperbarui",
			"status": 1,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "Order tidak ditemukan",
			"status": 0,
		})
	}
}

// DeleteOrder - Delete order by ID
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM orders WHERE id = ?"
	result, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	if affected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": "Order berhasil dihapus",
			"status": 1,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "Order tidak ditemukan",
			"status": 0,
		})
	}
}
