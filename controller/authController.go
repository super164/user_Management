package controller

import (
	"html/template"
	"net/http"
	"userManagement/service"
	"userManagement/session"
)

// Register 注册
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/register.html")
		t.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		confimPassword := r.FormValue("confirm_password")
		if username == "" || password == "" {
			w.Write([]byte("用户名或密码不能为空: "))
			return
		}
		if password != confimPassword {
			w.Write([]byte("两次密码不一致: "))
			return
		}
		err := service.Register(username, password)
		if err != nil {
			w.Write([]byte("注册失败: " + err.Error()))
			return
		}

		http.Redirect(w, r, "/login", 302)
	}
}

// Login 登录
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(w, nil)
		return
	}

	if r.Method == "POST" {

		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := service.Login(username, password)
		if err != nil {
			// 返回一段 JS 脚本，弹窗后返回登录页
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(err.Error()))
			return
		}

		// 设置 session
		session.CreateSession(w, *user)

		// 登录成功跳转用户列表
		http.Redirect(w, r, "/users", 302)
	}
}

// Logout 退出登录
func Logout(w http.ResponseWriter, r *http.Request) {
	session.DestroySession(w, r)
	http.Redirect(w, r, "/login", 302)
}
