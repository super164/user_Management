package router

import (
	"net/http" // 引入中间件包
	"userManagement/authMiddleware"
	"userManagement/controller"
)

func Init_router() {
	// 静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	http.HandleFunc("/", controller.InitHandler)
	http.HandleFunc("/index", controller.IndexHandler)
	http.HandleFunc("/register", controller.Register)
	http.HandleFunc("/login", controller.Login)

	// 使用 AuthMiddleware 包装需要鉴权的路由
	http.HandleFunc("/users", authMiddleware.AuthMiddleware(controller.ListUsers))
	http.HandleFunc("/logout", controller.Logout)
	http.HandleFunc("/users/delete", authMiddleware.AuthMiddleware(controller.DeleteUser))
	http.HandleFunc("/users/upload", authMiddleware.AuthMiddleware(controller.UploadAvatar))
	http.HandleFunc("/users/update", authMiddleware.AuthMiddleware(controller.UpdateUser))
	http.HandleFunc("/users/create", authMiddleware.AuthMiddleware(controller.CreateUser))

}
