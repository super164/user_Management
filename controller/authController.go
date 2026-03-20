package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"userManagement/service"
	"userManagement/session"
)

type jsonResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// InitHandler 登录跳转
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

// IndexHandler 跳转首页概括
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// 检查是否登录
	user, ok := session.GetSession(r)
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
	data := map[string]interface{}{
		"CurrentUser": user,
	}
	t.Execute(w, data)
}

func writeJSON(w http.ResponseWriter, status int, resp jsonResp) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

// Register 注册
func Register(w http.ResponseWriter, r *http.Request) {
	if _, ok := session.GetSession(r); ok {
		http.Redirect(w, r, "/users", http.StatusFound)
		return
	}

	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/register.html")
		_ = t.Execute(w, nil)
		return
	}

	if r.Method != "POST" {
		writeJSON(w, http.StatusMethodNotAllowed, jsonResp{Success: false, Message: "仅支持 POST 请求"})
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if username == "" || password == "" {
		writeJSON(w, http.StatusBadRequest, jsonResp{Success: false, Message: "用户名或密码不能为空"})
		return
	}
	if password != confirmPassword {
		writeJSON(w, http.StatusBadRequest, jsonResp{Success: false, Message: "两次密码不一致"})
		return
	}

	if err := service.Register(username, password); err != nil {
		writeJSON(w, http.StatusBadRequest, jsonResp{Success: false, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, jsonResp{Success: true, Message: "注册成功"})
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
