package dao

import (
	"errors"
	"userManagement/db"
	"userManagement/model"
)

// CreateUser 创建用户
func CreateUser(user model.User) error {
	_, err := db.DB.Exec(
		"INSERT INTO users(username, password,role)VALUES(?,?,?)",
		user.Username,
		user.Password,
		user.Role,
	)
	return err
}

// GetUserCount 获取用户总数
func GetUserCount(query string, status int) (int, error) {
	var count int
	// 查询总记录数
	sqlStr := "SELECT count(*) FROM users WHERE 1=1"
	var args []interface{}

	if query != "" {
		sqlStr += " AND username LIKE ?"
		args = append(args, "%"+query+"%")
	}
	// 状态筛选逻辑
	if status != -1 {
		sqlStr += " AND status = ?"
		args = append(args, status)
	}
	err := db.DB.QueryRow(sqlStr, args...).Scan(&count)
	return count, err
}

// AddUser 添加新用户
func AddUser(u model.User) error {
	//判断用户是否存在
	var count int
	err := db.DB.QueryRow("SELECT count(*) FROM users WHERE username=?", u.Username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}
	// 插入数据 (默认状态为 1-启用)
	_, err = db.DB.Exec("INSERT INTO users (username, password, role, avatar, status) VALUES (?, ?, ?, ?, ?)",
		u.Username, u.Password, u.Role, u.Avatar, 1)
	return err
}

// DeleteUser 删除用户
func DeleteUser(id int) error {
	_, err := db.DB.Exec("DELETE FROM users WHERE id=?", id)
	return err
}

// UpdateAvatar 更新头像
func UpdateAvatar(id int, path string) error {
	_, err := db.DB.Exec("UPDATE users SET  avatar=? WHERE id=?", path, id)
	return err
}

// UpdateUserDynamic 动态更新用户信息
// status = -1 代表不更新状态
// password = "" 代表不更新密码
// avatar = "" 代表不更新头像
func UpdateUserDynamic(id int, username, password string, status int, avatar string) error {
	// 1. 构建 SQL 语句
	query := "UPDATE users SET username=?"
	var args []interface{}
	args = append(args, username)

	// 2. 动态追加字段
	if password != "" {
		query += ", password=?"
		args = append(args, password)
	}

	if status != -1 {
		query += ", status=?" // 假设数据库里有 status 字段
		args = append(args, status)
	}

	if avatar != "" {
		query += ", avatar=?"
		args = append(args, avatar)
	}

	// 3. 添加 WHERE 条件
	query += " WHERE id=?"

	args = append(args, id)

	// 4. 执行
	_, err := db.DB.Exec(query, args...)
	return err
}

// UpdateLastLogin 更新最后登录时间
func UpdateLastLogin(id int, lastLogin string) error {
	_, err := db.DB.Exec("UPDATE users SET last_login=? WHERE id=?", lastLogin, id)
	return err
}
