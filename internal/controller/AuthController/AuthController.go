package Authcontroller

import (
	"encoding/json"
	"net/http"
	"strings"
	"ten_module/internal/DTO/request"
	"ten_module/internal/service/authservice"

	"github.com/gorilla/mux"
)

type AuthController struct {
	AuthService *authservice.AuthService
}

var AuthControllers *AuthController

func AuthControllerInit() {
	AuthControllers = &AuthController{
		AuthService: authservice.AuthServ,
	}
}
func (Controll *AuthController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/Login", Controll.Login).Methods("POST")
	r.HandleFunc("/refresh", Controll.RefreshToken).Methods("POST")
}
func (Controll *AuthController) Login(write http.ResponseWriter, req *http.Request) {
	var Login request.UserLogin
	errs := json.NewDecoder(req.Body).Decode(&Login)
	if errs != nil {
		http.Error(write, "Login Faile", http.StatusBadRequest)
		return
	}
	resp, err := Controll.AuthService.Login(Login)
	if err != nil {
		http.Error(write, "Login Faile", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(write).Encode(resp)
}
func (Controll *AuthController) RefreshToken(write http.ResponseWriter, req *http.Request) {
	Token := strings.Split(req.Header.Get("Authorization"), " ")[2]
	resp, err := Controll.AuthService.RefreshToken(Token)
	if err != nil {
		http.Error(write, "Refresh Failed", http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"Status": "OK",
		"Token":  resp,
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(write).Encode(response)
}
