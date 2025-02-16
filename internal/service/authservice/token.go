package authservice

import (
	"fmt"
	"log"
	config "ten_module/internal/Config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenHelper struct {
}
type Token interface {
	GenerateToken(username string, Role string) (string, error)
	VerifyToken(token string) string
	GetRoleToken(token string) string
}

var TokenUltils *TokenHelper

func GetTokenHelper() *TokenHelper {
	if TokenUltils != nil {
		TokenUltils = &TokenHelper{}
	}
	return TokenUltils
}
func Init() {
	TokenUltils = &TokenHelper{}
}

var Env = config.GetEnvConfig()

type TokenClaims struct {
	username string
	exp      time.Time
	Role     []string
	jwt.RegisteredClaims
}

func (h *TokenHelper) GenerateToken(username string, Role []string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"Role":     Role,
	})
	tokenString, err := token.SignedString([]byte(Env.JwtSecretKey()))
	if err != nil {
		log.Print(err)
		return "", err

	}
	return tokenString, nil
}
func (h *TokenHelper) VerifyToken(Token string) error {
	token, err := jwt.Parse(Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(Env.JwtSecretKey()), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
func (h *TokenHelper) GetRoleToken(Token string) ([]string, error) {
	claims := &TokenClaims{}
	_, err := jwt.ParseWithClaims(Token, claims, func(Token *jwt.Token) (interface{}, error) {
		return []byte(Env.JwtSecretKey()), nil
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return claims.Role, nil
}
