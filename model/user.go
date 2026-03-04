package model

type User struct {
	ID        int
	Username  string
	Password  string
	Role      string
	Avatar    string
	Status    int
	LastLogin string
}

//登录请求

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//注册请求

type RegisterResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
