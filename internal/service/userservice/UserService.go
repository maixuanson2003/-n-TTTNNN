package userservice

import (
	"errors"
	"log"
	constants "ten_module/internal/Constant"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repository.UserRepository
}
type MessageResponse struct {
	UserID  string
	Message string
	status  string
}
type Interface interface {
	UserRegister(UserReq request.UserRequest) (MessageResponse, error)
	GetListUser() ([]response.UserResponse, error)
	SearchUser(Name string, Age int, Email string, Address string, Role string, Gender string) ([]response.UserResponse, error)
	UpdateUser(UserReq request.UserRequest) (MessageResponse, error)
	DeleteUserById(Id int) (MessageResponse, error)
}

var UserServe *UserService

func GetUserService(UserRepo *repository.UserRepository) UserService {
	return UserService{
		UserRepo: UserRepo,
	}
}
func InitUserServ() {
	UserServe = &UserService{
		UserRepo: repository.UserRepo,
	}
}

func (service *UserService) UserRegister(UserReq request.UserRequest) (MessageResponse, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(UserReq.Password), 14)
	if err != nil {
		log.Print(" khong hash duoc mat khau")
		return MessageResponse{}, err
	}
	if UserReq.Username == "" ||
		UserReq.FullName == "" ||
		UserReq.Phone == "" || UserReq.Email == "" || UserReq.Address == "" || UserReq.Gender == "" {
		log.Print(errors.New("not enough data"))
		return MessageResponse{}, errors.New("not enough data")
	}
	User := entity.User{
		ID:       uuid.NewString(),
		Username: UserReq.Username,
		Password: string(bytes),
		FullName: UserReq.FullName,
		Phone:    UserReq.Phone,
		Email:    UserReq.Email,
		Address:  UserReq.Address,
		Gender:   UserReq.Gender,
		Age:      UserReq.Age,
		Role:     string(constants.USER),
	}
	repo := service.UserRepo
	errs := repo.Create(User)
	if errs != nil {
		return MessageResponse{}, err
	}
	return MessageResponse{
		status:  "SUCCESS",
		UserID:  User.ID,
		Message: "Ok",
	}, nil
}
func (service *UserService) GetListUser() ([]response.UserResponse, error) {
	repo := service.UserRepo
	ListUser, errors := repo.FindAll()
	if errors != nil {
		log.Print(errors)
		return nil, errors
	}
	var ListUserRes []response.UserResponse
	for _, user := range ListUser {
		Users := response.UserResponse{
			Phone:    user.Phone,
			FullName: user.FullName,
			Email:    user.Email,
			Address:  user.Address,
			Gender:   user.Gender,
			Age:      user.Age,
			Role:     user.Role,
		}
		ListUserRes = append(ListUserRes, Users)
	}
	return ListUserRes, nil
}
