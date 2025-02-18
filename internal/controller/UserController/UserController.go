package UserController

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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
	r.HandleFunc("/all", userController.GetListUser).Methods("GET")
	r.HandleFunc("/user/{id}", userController.DeleteUserById).Methods("DELETE")
	r.HandleFunc("/search", userController.SearchUser).Methods("POST")
	r.HandleFunc("/update/{id}", userController.UpdateUser).Methods("PUT")
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
func (userController *UserController) GetListUser(write http.ResponseWriter, Request *http.Request) {
	Resp, err := userController.UserService.GetListUser()
	if err != nil {
		http.Error(write, "Invalid request payload", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	json.NewEncoder(write).Encode(Resp)
}
func (userController *UserController) SearchUser(write http.ResponseWriter, Request *http.Request) {
	QueryParam := Request.URL.Query()
	// Get Queryparam
	Name := QueryParam.Get("fullname")
	Age := QueryParam.Get("age")
	var age int
	if Age == "" {
		age = 0
	} else {
		result, errs := strconv.Atoi(Age)
		if errs != nil {
			http.Error(write, "Invalid request payload", http.StatusBadRequest)
			return
		}
		age = result
	}
	Email := QueryParam.Get("email")
	Address := QueryParam.Get("address")
	Role := QueryParam.Get("role")
	Gender := QueryParam.Get("gender")
	//Get response call func SearchUser
	Resp, errs := userController.UserService.SearchUser(Name, age, Email, Address, Role, Gender)
	if errs != nil {
		http.Error(write, "failel To call Api", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	json.NewEncoder(write).Encode(Resp)
}
func (userController *UserController) DeleteUserById(write http.ResponseWriter, Request *http.Request) {
	url := Request.URL.Path
	userId := strings.Split(url, "/")[3]
	if userId == "" {
		http.Error(write, "failel To call Api", http.StatusBadRequest)
		return
	}
	resp, err := userController.UserService.DeleteUserById(userId)
	if err != nil {
		http.Error(write, "failel To call Api", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	json.NewEncoder(write).Encode(resp)
}
func (userController *UserController) UpdateUser(write http.ResponseWriter, Request *http.Request) {
	url := Request.URL.Path
	userId := strings.Split(url, "/")[3]
	if userId == "" {
		http.Error(write, "failel To call Api", http.StatusBadRequest)
		return
	}
	var UserRequest request.UserRequest
	err := json.NewDecoder(Request.Body).Decode(&UserRequest)
	if err != nil {
		http.Error(write, "failel To call Api", http.StatusBadRequest)
		return
	}
	resp, errs := userController.UserService.UpdateUser(UserRequest, userId)
	if errs != nil {
		http.Error(write, "failel To call Api", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	json.NewEncoder(write).Encode(resp)
}
