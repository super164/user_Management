package dao

import (
	"database/sql"
	"fmt"
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

// GetAllUsers 查询全部用户
func GetAllUsers() ([]model.User, error) {
	rows, err := db.DB.Query("SELECT id, username, password, role, avatar, status, IFNULL(last_login, '') FROM users")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.Avatar, &u.Status, &u.LastLogin)
		if err != nil {
			return nil, err
		}
		fmt.Println("读取到用户：", u.Username)
		users = append(users, u)
	}
	return users, nil
}

// DeleteUser 删除用户
func DeleteUser(id int) error {
	_, err := db.DB.Exec("DELETE FROM users WHERE id=?", id)
	return err
}

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

// GetUserCount 获取用户总数
func GetUserCount() (int, error) {
	var count int
	// 查询总记录数
	err := db.DB.QueryRow("SELECT count(*) FROM users").Scan(&count)
	return count, err
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
