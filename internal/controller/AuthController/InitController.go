package Authcontroller

import "ten_module/internal/service/authservice"

func InitController() {
	authservice.InitAuthService()
	AuthControllerInit()
}
