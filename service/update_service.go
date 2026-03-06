package service

import (
	"errors"
	"userManagement/dao"
	"userManagement/model"
)

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
