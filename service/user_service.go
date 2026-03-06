package service

import (
	"errors"
	"time"
	"userManagement/dao"
	"userManagement/model"
)

// Register 注册
func Register(username, password string) error {
	// 判断用户名是否存在
	existUser, err := dao.GetUserByName(username)
	if err != nil {
		return err
	}
	if existUser != nil {
		return errors.New("用户名已存在")
	}
	user := model.User{
		Username: username,
		Password: password,
		Role:     "user",
		Avatar:   "",
	}
	return dao.CreateUser(user)
}

// Login 登录验证
func Login(username, password string) (*model.User, error) {

	user, err := dao.GetUserByName(username)
	if err != nil {
		println("查询出错:", err.Error())
		// 如果查询出错（非 nil），直接返回错误，而不是 nil
		return nil, errors.New("查询失败: " + err.Error())
	}
	// 明确检查 user 是否为 nil
	if user == nil {
		println("数据库未查到用户")
		return nil, errors.New("用户名不存在")
	}
	// 用户存在后，再检查状态
	if user.Status == 0 {

		return nil, errors.New("该账号已被禁用，请联系管理员")
	}
	println("用户状态:", user.Status)
	println("数据库查到用户:", user.Username)
	// 最后检查密码
	if user.Password != password {
		return user, errors.New("密码错误")
	}
	now := time.Now().Format("2006-01-02 15:04:05")

	go func() {
		// 注意：UpdateLastLogin 第二个参数应该传时间字符串，或者修改 DAO 只传 ID
		// 这里假设您已经修改了 DAO 接收 username 或者 string
		_ = dao.UpdateLastLogin(user.ID, now)
	}()
	user.LastLogin = now
	return user, nil
}

// DeleteUser 删除用户
func DeleteUser(currentUser model.User, id int) error {
	if currentUser.Role != "admin" {
		return errors.New("无权限")
	}
	return dao.DeleteUser(id)
}
