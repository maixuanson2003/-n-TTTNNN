package authservice

import (
	"errors"
	"fmt"
	"log"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}
type Interface interface {
	Login(Request request.UserLogin) (string, error)
	RefreshToken(Token string) (string, error)
}

var AuthServ *AuthService

func InitAuthService() {
	fmt.Print(repository.UserRepo)
	AuthServ = &AuthService{
		UserRepo: repository.UserRepo,
	}
}

func (Service *AuthService) Login(Request request.UserLogin) (response.AuthResponse, error) {
	repo := Service.UserRepo
	UserArray, errs := repo.FindAll()
	if errs != nil {
		log.Print(errs)
		return response.AuthResponse{}, errs
	}
	UserFind := FindUserByUserName(Request.Username, UserArray)
	if UserFind == nil {
		return response.AuthResponse{}, errors.New("not found")
	}
	err := bcrypt.CompareHashAndPassword([]byte(UserFind.Password), []byte(Request.Password))
	if err != nil {
		return response.AuthResponse{}, errors.New("not match Password")
	}
	TokenHelper := TokenHelper{}
	tokenString, error := TokenHelper.GenerateToken(UserFind.Username, []string{UserFind.Role})
	if error != nil {
		return response.AuthResponse{}, errors.New("Gen token faild")
	}
	return response.AuthResponse{
		Username: UserFind.Username,
		Token:    tokenString,
		Role:     UserFind.Role,
		UserId:   UserFind.ID,
	}, nil
}
func FindUserByUserName(Username string, UserArray []entity.User) *entity.User {
	for _, User := range UserArray {
		if User.Username == Username {
			return &User
		}
	}
	return nil
}
func (serviec *AuthService) RefreshToken(Token string) (string, error) {
	err := TokenUltils.VerifyToken(Token)
	if err == nil {
		return "", errors.New("Token is valid")
	}
	claims := &TokenClaims{}
	_, errs := jwt.ParseWithClaims(Token, claims, func(Token *jwt.Token) (interface{}, error) {
		return []byte(Env.JwtSecretKey()), nil
	})
	if errs != nil {
		return "", errors.New("failde parse")
	}
	NewToken, errors := TokenUltils.GenerateToken(claims.username, claims.Role)
	if errors != nil {
		return "", errors
	}
	return NewToken, nil
}
