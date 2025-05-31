package userservice

import (
	"errors"
	"log"
	"regexp"
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
	OtpRepo  *repository.OtpRepository
}
type MessageResponse struct {
	UserID  string
	Message string
	Status  string
}
type Interface interface {
	UserRegister(UserReq request.UserRequest) (MessageResponse, error)
	GetListUser() ([]response.UserResponse, error)
	GetUserById(Id string) (response.UserResponse, error)
	SearchUser(Name string, Age int, Email string, Address string, Role string, Gender string) ([]response.UserResponse, error)
	UpdateUser(UserReq request.UserRequest, Id string) (MessageResponse, error)
	DeleteUserById(Id string) (MessageResponse, error)
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
		OtpRepo:  repository.OtpRepo,
	}
}
func UserReqMapToUserEntity(UserReq request.UserRequest, Byte []byte) entity.User {
	return entity.User{
		ID:       uuid.NewString(),
		Username: UserReq.Username,
		Password: string(Byte),
		FullName: UserReq.FullName,
		Phone:    UserReq.Phone,
		Email:    UserReq.Email,
		Address:  UserReq.Address,
		Gender:   UserReq.Gender,
		Age:      UserReq.Age,
		Role:     string(constants.ADMIN),
	}
}
func UserEntityMapToUserResponse(user entity.User) response.UserResponse {
	return response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Phone:    user.Phone,
		FullName: user.FullName,
		Email:    user.Email,
		Address:  user.Address,
		Gender:   user.Gender,
		Age:      user.Age,
		Role:     user.Role,
	}
}

func (service *UserService) UserRegister(UserReq request.UserRequest, types string) (MessageResponse, error) {
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
	User := UserReqMapToUserEntity(UserReq, bytes)
	if types == "user" {
		User.Role = string(constants.USER)
	}

	repo := service.UserRepo

	errs := repo.Create(User)
	if errs != nil {
		return MessageResponse{}, err
	}
	return MessageResponse{
		Status:  "SUCCESS",
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
		Users := UserEntityMapToUserResponse(user)
		ListUserRes = append(ListUserRes, Users)
	}
	return ListUserRes, nil
}
func (service *UserService) SearchUser(Name string, Age int, Email string, Address string, Role string, Gender string) ([]response.UserResponse, error) {
	repo := service.UserRepo
	ListUser, err := repo.GetUserQuery(Name, Age, Email, Address, Role, Gender)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var ListUserRes []response.UserResponse
	for _, user := range ListUser {
		Users := UserEntityMapToUserResponse(user)
		ListUserRes = append(ListUserRes, Users)
	}
	return ListUserRes, nil
}
func (service *UserService) DeleteUserById(Id string) (MessageResponse, error) {
	repo := service.UserRepo
	err := repo.DeleteById(Id)
	if err != nil {
		return MessageResponse{
			Status:  "Fail",
			UserID:  Id,
			Message: "Ok",
		}, err

	}
	return MessageResponse{
		Status:  "SUCCESS",
		UserID:  Id,
		Message: "Ok",
	}, nil
}
func (service *UserService) UpdateUser(UserReq request.UserUpdate, Id string) (MessageResponse, error) {
	repo := service.UserRepo
	User, checkFaild := repo.FindById(Id)
	if checkFaild != nil {
		return MessageResponse{
			Status:  "FAILED",
			UserID:  "s",
			Message: "bad request",
		}, checkFaild
	}

	Validate := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !Validate.MatchString(UserReq.Email) {
		return MessageResponse{
			Status:  "FAILED",
			UserID:  "s",
			Message: "not valid email",
		}, errors.New("not valid email")
	}
	if len(UserReq.Phone) < 10 {
		return MessageResponse{
			Status:  "FAILED",
			UserID:  "s",
			Message: "not valid Phone Number",
		}, errors.New("not valid Phone Number")
	}
	User.Username = UserReq.Username
	User.Phone = UserReq.Phone
	User.Address = UserReq.Address
	User.Age = UserReq.Age
	User.Email = UserReq.Email
	User.FullName = UserReq.FullName
	User.Gender = UserReq.Gender
	IsFaildToUpdate := repo.Update(User, Id)
	if IsFaildToUpdate != nil {
		return MessageResponse{
			Status:  "FAILED",
			UserID:  Id,
			Message: "Update Faild",
		}, IsFaildToUpdate

	}
	return MessageResponse{
		Status:  "Success",
		UserID:  Id,
		Message: "Update Success",
	}, nil
}
func (service *UserService) GetUserById(Id string) (response.UserResponse, error) {
	repo := service.UserRepo
	User, checkFaild := repo.FindById(Id)
	if checkFaild != nil {
		return response.UserResponse{}, nil
	}
	UserResponse := UserEntityMapToUserResponse(User)
	return UserResponse, nil
}
