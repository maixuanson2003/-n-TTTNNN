package UserController

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	middleware "ten_module/Middleware"
	"ten_module/internal/DTO/request"
	"ten_module/internal/repository"
	"ten_module/internal/service/userservice"
	"time"

	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	UserService *userservice.UserService
	MiddleWare  *middleware.UseMiddleware
}

var UserControll *UserController

func UserControllerInit() {
	UserControll = &UserController{
		UserService: userservice.UserServe,
		MiddleWare:  middleware.Middlewares,
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
	middleware := userController.MiddleWare
	r.HandleFunc("/register", userController.UserRegister).Methods("POST")
	r.HandleFunc("/createuser", middleware.Chain(userController.UserCreate, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("POST")
	r.HandleFunc("/all", userController.GetListUser).Methods("GET")
	r.HandleFunc("/user/{id}", middleware.Chain(userController.DeleteUserById, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("DELETE")
	r.HandleFunc("/search", middleware.Chain(userController.SearchUser, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN", "USER"}))).Methods("POST")
	// r.HandleFunc("/update/{id}", middleware.Chain(userController.UpdateUser, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("PUT")
	r.HandleFunc("/update/{id}", middleware.Chain(userController.UpdateUser, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("PUT")
	r.HandleFunc("/getuser/{id}", userController.GetUserById).Methods("GET")
	r.HandleFunc("/deleteuser/{id}", userController.DeleteUserById).Methods("DELETE")
	r.HandleFunc("/admin/users/export", middleware.Chain(userController.ExportUsersExcel, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("GET")
	r.HandleFunc("/change/password", middleware.Chain(userController.ChangePassword, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN", "USER"})))
}
func (userController *UserController) UserRegister(write http.ResponseWriter, Request *http.Request) {
	var Body request.UserRequest
	err := json.NewDecoder(Request.Body).Decode(&Body)
	if err != nil {
		http.Error(write, "Invalid request payload", http.StatusBadRequest)
		return
	}

	Resp, err := userController.UserService.UserRegister(Body, "user")
	if err != nil {
		http.Error(write, "Invalid request payload", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	json.NewEncoder(write).Encode(Resp)

}
func (userController *UserController) ExportUsersExcel(w http.ResponseWriter, r *http.Request) {
	users, err := userController.UserService.UserRepo.FindAll() // Hàm lấy tất cả user từ DB
	if err != nil {
		http.Error(w, "Không thể lấy dữ liệu người dùng", http.StatusInternalServerError)
		return
	}

	f := excelize.NewFile()
	sheet := "Users"
	f.NewSheet(sheet)

	headers := []string{"ID", "Username", "Email", "Phone", "Gender", "Age", "Role"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	for idx, u := range users {
		row := idx + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), u.ID)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), u.Username)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), u.Email)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), u.Phone)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), u.Gender)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), u.Age)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), u.Role)
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=users_"+time.Now().Format("20060102")+".xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")

	if err := f.Write(w); err != nil {
		http.Error(w, "Lỗi khi tạo file Excel", http.StatusInternalServerError)
	}
}
func (userController *UserController) UserCreate(write http.ResponseWriter, Request *http.Request) {
	var Body request.UserRequest
	err := json.NewDecoder(Request.Body).Decode(&Body)
	if err != nil {
		http.Error(write, "Invalid request payload", http.StatusBadRequest)
		return
	}
	Resp, err := userController.UserService.UserRegister(Body, "admin")
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
	var UserRequest request.UserUpdate
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
func (userController *UserController) GetUserById(write http.ResponseWriter, Request *http.Request) {
	url := Request.URL.Path
	userId := strings.Split(url, "/")[3]
	if userId == "" {
		http.Error(write, "failel To call Api", http.StatusBadRequest)
		return
	}
	resp, err := userController.UserService.GetUserById(userId)
	if err != nil {
		http.Error(write, "failel To call Api", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	json.NewEncoder(write).Encode(resp)
}
func (userController *UserController) ChangePassword(write http.ResponseWriter, Request *http.Request) {
	userid := Request.URL.Query().Get("userid")
	newPassword := Request.URL.Query().Get("newpassword")
	repo := userController.UserService.UserRepo
	user, err := repo.FindById(userid)
	if err != nil {
		http.Error(write, "not found", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	user.Password = string(hashedPassword)
	errs := repo.Update(user, userid)
	if errs != nil {
		http.Error(write, "update failed", http.StatusBadRequest)
		return
	}

	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	json.NewEncoder(write).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Đổi mật khẩu thành công",
	})
}
