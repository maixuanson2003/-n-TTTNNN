package server

import (
	database "ten_module/Database"
	middleware "ten_module/Middleware"
	Authcontroller "ten_module/internal/controller/AuthController"
	"ten_module/internal/controller/UserController"
	"ten_module/internal/repository"
	"ten_module/internal/service/authservice"
)

func InitSingleton() {
	middleware.InitMiddleWare()
	authservice.Init()
	database.Init()
	repository.InitUserRepo()
	UserController.InitService()
	Authcontroller.InitController()
}
