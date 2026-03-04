package service

import (
	"errors"
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
	go func() {
		// 注意：UpdateLastLogin 第二个参数应该传时间字符串，或者修改 DAO 只传 ID
		// 这里假设您已经修改了 DAO 接收 username 或者 string
		_ = dao.UpdateLastLogin(user.ID, username)
	}()
	return user, nil
}

// DeleteUser 删除用户
func DeleteUser(currentUser model.User, id int) error {
	if currentUser.Role != "admin" {
		return errors.New("无权限")
	}
	return dao.DeleteUser(id)
}

// UpdateAvatar 更新头像
func UpdateAvatar(currentUser model.User, targetId int, path string) error {
	if currentUser.Role == "admin" {
		return dao.UpdateAvatar(targetId, path)
	}
	if currentUser.ID == targetId {
		return dao.UpdateAvatar(targetId, path)
	}
	return errors.New("无权限")
}

// UpdateUserInfo 综合更新用户信息
func UpdateUserInfo(currentUser model.User, targetID int, username, password string, status int, avatar string) error {
	// 1. 权限检查
	if currentUser.Role != "admin" {
		// 如果不是管理员，必须是修改自己的信息
		if currentUser.ID != targetID {
			return errors.New("无权修改他人信息")
		}
		// 普通用户不能修改状态，强制设为 -1 (代表不修改) 或者保持原样
		// 这里我们在 DAO 层处理：如果 status 为 -1 则不更新该字段
		status = -1
	}
	// 如果尝试修改状态为禁用 (status == 0)
	if status == 0 {
		// 先查询目标用户的信息，确认其角色
		targetUser, err := dao.GetUserByID(targetID)
		if err != nil {
			return err
		}
		// 如果目标用户是管理员，禁止禁用
		if targetUser.Role == "admin" {
			return errors.New("无法禁用管理员账号")
		}
	}
	// 2. 调用 DAO 更新
	return dao.UpdateUserDynamic(targetID, username, password, status, avatar)
}
