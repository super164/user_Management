package main

import (
	"log"
	"net/http"
	"userManagement/controller"
	"userManagement/db"

	_ "github.com/go-sql-driver/mysql"
)

// 测试首页
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("用户管理系统启动成功"))

	w.Write([]byte("插入成功"))
}
func main() {

	//初始化数据库
	db.InitDB()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/register", controller.Register)
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/users", controller.ListUsers)
	http.HandleFunc("/logout", controller.Logout)
	http.HandleFunc("/users/delete", controller.DeleteUser)
	http.HandleFunc("/users/update", controller.UpdateUser)

	log.Println("服务器启动: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("启动失败:", err)
	}
}
