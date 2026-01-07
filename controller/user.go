package controller

import (
	"crud-go/config"
	"crud-go/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUser(c *gin.Context) {
	var users []model.User

	// Get pagination parameters from query
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "5")

	// Convert page and limit to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 5
	}

	offset := (pageInt - 1) * limitInt

	// Query to count total users
	var total int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": err.Error(),
			"status": 0,
		})
		return
	}

	// Query with LIMIT and OFFSET for pagination
	rows, err := config.DB.Query("SELECT * FROM users LIMIT ? OFFSET ?", limitInt, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": err.Error(),
			"status": 0,
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.Id,
			&u.Nama,
			&u.Email,
			&u.Password,
			&u.Role,
			&u.Created_at,
			&u.Updated_at,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": err.Error(),
				"status": 0,
			})
			return
		}
		users = append(users, u)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": users,
		"status": 1,
		"page":   pageInt,
		"limit":  limitInt,
		"total":  total,
	})
}

func SetUser(c *gin.Context) {
	var newUser model.User

	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	// Set default role to customer if not provided
	if newUser.Role == "" {
		newUser.Role = "customer"
	}

	// Validate role
	validRoles := map[string]bool{"admin": true, "customer": true, "seller": true}
	if !validRoles[newUser.Role] {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    "Invalid role. Must be admin, customer, or seller",
			"status": http.StatusBadRequest,
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    "Error hashing password",
			"status": http.StatusBadRequest,
		})
		return
	}

	// Set createdAt and updatedAt to current time
	newUser.Created_at = time.Now()
	newUser.Updated_at = time.Now()

	var u = "INSERT INTO users (nama, email, password, role, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?)"
	rows, err := config.DB.Exec(u, newUser.Nama, newUser.Email, string(hashedPassword), newUser.Role, newUser.Created_at, newUser.Updated_at)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	if affected > 0 {
		newUser.Password = ""
		c.JSON(http.StatusOK, gin.H{
			"result": "Sukses Menambahkan User",
			"status": http.StatusOK,
			"data":   newUser,
		})
	}
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser model.User

	err := c.BindJSON(&updatedUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}
	updatedUser.Updated_at = time.Now()

	// Validate role if provided
	if updatedUser.Role != "" {
		validRoles := map[string]bool{"admin": true, "customer": true, "seller": true}
		if !validRoles[updatedUser.Role] {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":    "Invalid role. Must be admin, customer, or seller",
				"status": http.StatusBadRequest,
			})
			return
		}
	}

	query := "UPDATE users SET nama = ?, email = ?, password = ?, role = ?, updatedAt = ? WHERE id = ?"
	result, err := config.DB.Exec(query, updatedUser.Nama, updatedUser.Email, updatedUser.Password, updatedUser.Role, updatedUser.Updated_at, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	if affected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": "User updated successfully",
			"status": http.StatusOK,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "User not found",
			"status": http.StatusNotFound,
		})
	}
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM users WHERE id = ?"
	result, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	if affected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": "User deleted successfully",
			"status": http.StatusOK,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "User not found",
			"status": http.StatusNotFound,
		})
	}
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user model.User

	query := "SELECT id, nama, email, password, role, createdAt, updatedAt FROM users WHERE id = ?"
	err := config.DB.QueryRow(query, id).Scan(
		&user.Id,
		&user.Nama,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "User not found",
			"status": http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": user,
		"status": http.StatusOK,
	})
}

// Login - User login dengan email dan password
func Login(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.BindJSON(&loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": err.Error(),
			"status": 0,
		})
		return
	}

	var user model.User

	// Query user by email
	query := "SELECT id, nama, email, password, role, createdAt, updatedAt FROM users WHERE email = ?"
	err = config.DB.QueryRow(query, loginRequest.Email).Scan(
		&user.Id,
		&user.Nama,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": "Email atau password salah",
			"status": 0,
		})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": "Email atau password salah",
			"status": 0,
		})
		return
	}

	// Return user data without password
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"result": user,
		"status": 1,
		"message": "Login berhasil",
	})
}

// Register - User registration dengan hashing password
func Register(c *gin.Context) {
	var newUser model.User

	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	// Set default role to customer if not provided
	if newUser.Role == "" {
		newUser.Role = "customer"
	}

	// Validate role
	validRoles := map[string]bool{"admin": true, "customer": true, "seller": true}
	if !validRoles[newUser.Role] {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    "Invalid role. Must be admin, customer, or seller",
			"status": 0,
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    "Error hashing password",
			"status": 0,
		})
		return
	}

	newUser.Created_at = time.Now()
	newUser.Updated_at = time.Now()

	query := "INSERT INTO users (nama, email, password, role, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := config.DB.Exec(query, newUser.Nama, newUser.Email, string(hashedPassword), newUser.Role, newUser.Created_at, newUser.Updated_at)
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
		newUser.Password = ""
		c.JSON(http.StatusOK, gin.H{
			"result":  "User berhasil didaftarkan",
			"status":  1,
			"data":    newUser,
		})
	}
}