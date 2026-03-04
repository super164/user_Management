package controller

import (
	"fmt"
	"html/template"
	"io"

	"net/http"
	"os"
	"strconv"
	"time"
	"userManagement/dao"

	"userManagement/service"
	"userManagement/session"
)

// ListUsers 用户主页
func ListUsers(w http.ResponseWriter, r *http.Request) {
	// 检查是否登录
	currentUser, ok := session.GetSession(r)
	if !ok {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// 查询所有用户
	users, err := dao.GetAllUsers()
	if err != nil {
		http.Error(w, "查询用户失败"+err.Error(), http.StatusInternalServerError)
		return
	}
	//解析模版
	t, err := template.ParseFiles("templates/users.html")
	if err != nil {
		http.Error(w, "模版解析失败"+err.Error(), http.StatusInternalServerError)
		return
	}
	// 传递数据给模板
	data := map[string]interface{}{
		"CurrentUser": currentUser,
		"Users":       users,
	}
	//渲染页面
	t.Execute(w, data)
}

// DeleteUser 删除用户
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := session.GetSession(r)
	if !ok {
		http.Redirect(w, r, "/login", 302)
		return
	}
	// 检查权限
	if currentUser.Role != "admin" {
		w.Write([]byte("无权限操作"))
		return
	}
	//获取ID
	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)

	//删除
	err := service.DeleteUser(currentUser, id)
	if err != nil {
		w.Write([]byte("删除失败: " + err.Error()))
		return
	}
	http.Redirect(w, r, "/users", 302)
}

// UploadAvatar 上传头像
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// 必须是 POST
	if r.Method != "POST" {
		http.Error(w, "非法请求", http.StatusMethodNotAllowed)
		return
	}

	//检查是否登录
	currentUser, ok := session.GetSession(r)
	if !ok {
		//未登录则跳转到登录页
		http.Redirect(w, r, "/login", 302)
		return
	}

	// 解析表单
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "解析失败", 500)
		return
	}

	// 4. 获取目标用户 ID 并转换类型
	targetIDStr := r.FormValue("user_id")
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "文件获取失败", 500)
		return
	}
	defer file.Close()

	// 生成文件名（避免覆盖）
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)

	_ = os.MkdirAll("./uploads", 0755)

	// 创建文件
	dst, err := os.Create("./uploads/" + filename)
	if err != nil {
		http.Error(w, "文件保存失败", 500)
		return
	}
	defer dst.Close()

	// 写入文件内容
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "文件写入失败", http.StatusInternalServerError)
		return
	}
	// 7.调用 Service 层进行权限验证和更新
	// 传入 currentUser 用于鉴权，targetID 用于指定要修改的用户
	err = service.UpdateAvatar(currentUser, targetID, filename)
	if err != nil {
		// 如果是权限错误或数据库错误
		http.Error(w, "更新失败: "+err.Error(), http.StatusForbidden)
		// 如果更新数据库失败，最好把刚才上传的垃圾文件删掉（可选优化）
		_ = os.Remove("./uploads/" + filename)
		return
	}

	// 8. 成功后重定向回用户列表
	http.Redirect(w, r, "/users", 302)
}

// UpdateUser 更新用户信息
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// 1. 必须是 POST 请求
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. 检查登录
	currentUser, ok := session.GetSession(r)
	if !ok {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// 3. 解析表单 (支持文件上传)
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "解析表单失败", http.StatusBadRequest)
		return
	}

	// 4. 获取参数
	targetIDStr := r.FormValue("id")
	username := r.FormValue("username")
	password := r.FormValue("password") // 如果为空则不修改
	statusStr := r.FormValue("status")

	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	// 状态默认为 1 (启用)，如果前端传了则用前端的
	status := 1
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}

	// 5. 处理头像上传
	var avatarPath string
	file, handler, err := r.FormFile("avatar")
	if err == nil {
		// 有文件上传
		defer file.Close()
		// 生成唯一文件名
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)

		// 确保目录存在
		_ = os.MkdirAll("./uploads", 0755)

		// 保存文件
		dst, err := os.Create("./uploads/" + filename)
		if err != nil {
			http.Error(w, "文件保存失败", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		io.Copy(dst, file)
		avatarPath = filename
	}

	// 6. 调用 Service 层进行更新
	// 参数说明: 当前操作人, 目标用户ID, 新用户名, 新密码, 新状态, 新头像路径
	err = service.UpdateUserInfo(currentUser, targetID, username, password, status, avatarPath)
	if err != nil {
		w.WriteHeader(http.StatusForbidden) // 或者 200，前端根据 success 字段判断
		w.Write([]byte(fmt.Sprintf(`{"success": false, "message": "更新失败: %s"}`, err.Error())))
		return
	}
	//7. 同步更新 Session (如果是修改当前登录用户)
	if currentUser.ID == targetID {

		// 更新变动的字段
		if username != "" {
			currentUser.Username = username
		}
		if avatarPath != "" {
			currentUser.Avatar = avatarPath
		}
		// 如果改了密码，理论上也要更新
		if password != "" {
			currentUser.Password = password
		}

		// 更新 Session
		session.UpdateSessionUser(w, r, &currentUser)
	}
	// 8. 成功后返回 JSON 成功信息
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "message": "更新成功"}`))
}
