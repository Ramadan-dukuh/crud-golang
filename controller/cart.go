package controller

import (
	"crud-go/config"
	"crud-go/model"
	"net/http"
	_ "strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetCartByUserID - Get cart items for a specific user
func GetCartByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	var carts []model.Cart

	query := "SELECT id, user_id, barang_id, qty, createdAt FROM cart WHERE user_id = ?"
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
		var cart model.Cart
		err := rows.Scan(
			&cart.Id,
			&cart.User_id,
			&cart.Barang_id,
			&cart.Qty,
			&cart.Created_at,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": err.Error(),
				"status": 0,
			})
			return
		}
		carts = append(carts, cart)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": carts,
		"status": 1,
	})
}

// AddToCart - Add item to cart (Customer only)
func AddToCart(c *gin.Context) {
	// Check role - only customer can add to cart
	userRole := GetUserRoleFromContext(c)
	if userRole != "customer" {
		c.JSON(http.StatusForbidden, gin.H{
			"err":    "Hanya customer yang dapat menambahkan ke cart",
			"status": 0,
		})
		return
	}

	var newCart model.Cart

	err := c.BindJSON(&newCart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	newCart.Created_at = time.Now()

	query := "INSERT INTO cart (user_id, barang_id, qty, createdAt) VALUES (?, ?, ?, ?)"
	result, err := config.DB.Exec(query, newCart.User_id, newCart.Barang_id, newCart.Qty, newCart.Created_at)
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
			"result": "Item berhasil ditambahkan ke cart",
			"status": 1,
			"data":   newCart,
		})
	}
}

// UpdateCart - Update cart item quantity
func UpdateCart(c *gin.Context) {
	cartID := c.Param("cart_id")
	var updatedCart model.Cart

	err := c.BindJSON(&updatedCart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	query := "UPDATE cart SET qty = ? WHERE id = ?"
	result, err := config.DB.Exec(query, updatedCart.Qty, cartID)
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
			"result": "Cart berhasil diperbarui",
			"status": 1,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "Cart tidak ditemukan",
			"status": 0,
		})
	}
}

// RemoveFromCart - Remove item from cart
func RemoveFromCart(c *gin.Context) {
	cartID := c.Param("cart_id")

	query := "DELETE FROM cart WHERE id = ?"
	result, err := config.DB.Exec(query, cartID)
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
			"result": "Item berhasil dihapus dari cart",
			"status": 1,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "Cart tidak ditemukan",
			"status": 0,
		})
	}
}

// ClearCart - Clear all cart items for a user
func ClearCart(c *gin.Context) {
	userID := c.Param("user_id")

	query := "DELETE FROM cart WHERE user_id = ?"
	result, err := config.DB.Exec(query, userID)
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
			"result": "Cart berhasil dikosongkan",
			"status": 1,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "Cart tidak ditemukan",
			"status": 0,
		})
	}
}
