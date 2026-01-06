package controller

import (
	"crud-go/config"
	"crud-go/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllUser(c *gin.Context) {
	var user []model.User

	rows, err := config.DB.Query("select * from users")
	if err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"result" : err.Error(),
			"status" : 0,
		})
		return
	}
	defer rows.Close()

	for rows.Next(){
		var u model.User
		err := rows.Scan(
			&u.Id,
			&u.Nama,
			&u.Email,
			&u.Created_at,
			&u.Updated_at,
		)
		if err != nil {
			fmt.Println("Error scanning users")
			c.JSON(http.StatusInternalServerError,gin.H{
				"result" : err,
				"status" : 0,
			})
			panic(err)
		}
		fmt.Print(u)
		user = append(user,u)
	}
	if len(user) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result" : user,
			"status" : 0,
		})
	}
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

    // Set createdAt and updatedAt to current time
    newUser.Created_at = time.Now()
    newUser.Updated_at = time.Now()

    var u = "INSERT INTO users (nama, email, createdAt, updatedAt) VALUES (?, ?, ?, ?)"
    rows, err := config.DB.Exec(u, newUser.Nama, newUser.Email, newUser.Created_at, newUser.Updated_at)
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

	query := "UPDATE users SET nama = ?, email = ?, updatedAt = ? WHERE id = ?"
	result, err := config.DB.Exec(query, updatedUser.Nama, updatedUser.Email, updatedUser.Updated_at, id)
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

	query := "SELECT id, nama, email, createdAt, updatedAt FROM users WHERE id = ?"
	err := config.DB.QueryRow(query, id).Scan(
		&user.Id,
		&user.Nama,
		&user.Email,
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