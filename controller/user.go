package controller

import (
	"crud-go/config"
	"crud-go/model"
	"fmt"
	"net/http"

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
func SetUser(c*gin.Context){
	var newUser model.User

	err := c.BindJSON(&newUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err" : err,
			"status" : http.StatusBadRequest,
		})
	}
	
	var u = "insert into users (nama,email,createdAt,updatedAt) values (?,?,?,?)"
	rows, err := config.DB.Exec(u,newUser.Nama,newUser.Email,newUser.Created_at,newUser.Updated_at)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err" : err,
			"status" : http.StatusBadRequest,
		})
	}
	affected, err := rows.RowsAffected()
	if err != err {
		c.JSON(http.StatusBadRequest, gin.H{
			"err" : err,
			"status" : http.StatusBadRequest,
		})
	}
	if affected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result" : "Sukses Menambahkan User",
			"status" : http.StatusOK,
			"data" : newUser,
		})
	}
}