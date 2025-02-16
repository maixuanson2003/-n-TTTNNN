package UserController

import (
	"encoding/json"
	"net/http"
	"ten_module/internal/DTO/request"
	"ten_module/internal/repository"
	"ten_module/internal/service/userservice"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserController struct {
	UserService *userservice.UserService
}

var UserControll *UserController

func UserControllerInit() {
	UserControll = &UserController{
		UserService: userservice.UserServe,
	}
}

func GetUserController(Database *gorm.DB) UserController {
	Repo := repository.UserRepository{
		DB: Database,
	}
	UserService := userservice.GetUserService(&Repo)
	return UserController{
		UserService: &UserService,
	}
}
func (userController *UserController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/register", userController.UserRegister).Methods("POST")

}
func (userController *UserController) UserRegister(write http.ResponseWriter, Request *http.Request) {
	var Body request.UserRequest
	err := json.NewDecoder(Request.Body).Decode(&Body)
	if err != nil {
		http.Error(write, "Invalid request payload", http.StatusBadRequest)
		return
	}

	Resp, err := userController.UserService.UserRegister(Body)
	if err != nil {
		http.Error(write, "Invalid request payload", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	json.NewEncoder(write).Encode(Resp)

}
