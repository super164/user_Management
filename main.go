package main

import (
	"log"
	"net/http"
	"userManagement/db"
	"userManagement/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//初始化数据库
	db.InitDB()
	//注册路由
	router.Init_router()

	log.Println("服务器启动: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("启动失败:", err)
	}
}
