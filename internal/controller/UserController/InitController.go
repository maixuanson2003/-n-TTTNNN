package UserController

import (
	"ten_module/internal/service/userservice"
)

func InitService() {

	userservice.InitUserServ()
	UserControllerInit()
}
