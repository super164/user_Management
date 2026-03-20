package service

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var usernameRe = regexp.MustCompile(`^[A-Za-z0-9]{5,10}$`)

func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("用户名不能为空")
	}
	if strings.TrimSpace(username) != username {
		return errors.New("用户名不能包含空格")
	}
	if !usernameRe.MatchString(username) {
		return errors.New("用户名格式不合法：5-10位，仅允许字母和数字")
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}
	if len(password) < 6 || len(password) > 10 {
		return errors.New("密码长度不合法：6-10位")
	}
	for _, r := range password {
		if unicode.IsSpace(r) {
			return errors.New("密码不能包含空格")
		}
	}
	return nil
}
