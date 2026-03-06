package dao

import (
	"database/sql"
	"userManagement/db"
	"userManagement/model"
)

// GetUserByName 根据用户名查询
func GetUserByName(username string) (*model.User, error) {
	var user model.User
	err := db.DB.QueryRow(
		"SELECT id, username, password, role, avatar, status, IFNULL(last_login, '') FROM users WHERE username=?",
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Avatar,
		&user.Status,
		&user.LastLogin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &user, err
}

// GetUserByID 根据ID查询用户
func GetUserByID(id int) (*model.User, error) {
	var user model.User
	// 记得 Scan 所有字段，包括 status 和 last_login
	err := db.DB.QueryRow(
		"SELECT id, username, password, role, avatar, status, IFNULL(last_login, '') FROM users WHERE id=?",
		id,
	).Scan(
		&user.ID, &user.Username, &user.Password, &user.Role, &user.Avatar, &user.Status, &user.LastLogin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsersByPage 分页查询用户
func GetUsersByPage(page, pageSize int) ([]model.User, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 执行分页查询
	rows, err := db.DB.Query("SELECT id, username, password, role, avatar, status, IFNULL(last_login, '') FROM users LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		// 扫描数据（保持和 GetAllUsers 一致）
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.Avatar, &u.Status, &u.LastLogin)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
