package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	var db_detail = "root:@tcp(127.0.0.1:3306)/crud_dbnode?parseTime=true"

	dbConn, err := sql.Open("mysql",db_detail)
	if err != nil {
		fmt.Println("Database Connection failed")
		panic(err.Error())
	}else{
		fmt.Println("Database Connected")
	}

	DB = dbConn
}