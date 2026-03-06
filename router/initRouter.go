package router

import (
	"net/http"
	"userManagement/controller"
)

func Init_router() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	http.HandleFunc("/", controller.InitHandler)
	http.HandleFunc("/index", controller.IndexHandler)
	http.HandleFunc("/register", controller.Register)
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/users", controller.ListUsers)
	http.HandleFunc("/logout", controller.Logout)
	http.HandleFunc("/users/delete", controller.DeleteUser)
	http.HandleFunc("/users/update", controller.UpdateUser)
}
