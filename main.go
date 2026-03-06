package main

import (
	"html/template"
	"log"
	"net/http"
	"userManagement/controller"
	"userManagement/db"
	"userManagement/session"

	_ "github.com/go-sql-driver/mysql"
)

// 登录跳转
func InitHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}
	//检查是否已登录
	_, ok := session.GetSession(r)
	if ok {
		//已登录
		http.Redirect(w, r, "/users", http.StatusFound)
	} else {
		//未登录
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// 跳转首页概括
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// 检查是否登录
	_, ok := session.GetSession(r)
	if !ok {
		// 未登录则跳转到登录页
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// 渲染 index.html 模板
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "模板解析失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func main() {

	//初始化数据库
	db.InitDB()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	http.HandleFunc("/", InitHandler)
	http.HandleFunc("/index", IndexHandler)
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
