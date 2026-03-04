package db

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB() {

	dsn := "root:123456@tcp(127.0.0.1:3306)/user_management?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("数据库Ping失败:", err)
	}

	log.Println("数据库连接成功")
}
