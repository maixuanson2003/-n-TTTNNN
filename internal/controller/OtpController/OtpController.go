package otpcontroller

import (
	"encoding/json"
	"net/http"
	"ten_module/internal/service/otpservice"

	"github.com/gorilla/mux"
)

type OtpController struct {
	OtpServe *otpservice.OtpService
}

var OtpControll *OtpController

func InitOtpController() {
	OtpControll = &OtpController{
		OtpServe: otpservice.OtpServe,
	}
}
func (OtpControll *OtpController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/checkotp", OtpControll.CheckOtp).Methods("POST")
	r.HandleFunc("/sendotp", OtpControll.SendOtpHandler).Methods("POST")

}
func (OtpControll *OtpController) CheckOtp(write http.ResponseWriter, request *http.Request) {
	otp := request.URL.Query().Get("otp")
	check := OtpControll.OtpServe.CheckOtp(otp)

	write.Header().Set("Content-Type", "application/json")

	if check {
		response := map[string]interface{}{
			"success": true,
			"message": "OTP is valid",
		}
		json.NewEncoder(write).Encode(response)
	} else {
		write.WriteHeader(http.StatusUnauthorized)
		response := map[string]interface{}{
			"success": false,
			"message": "Invalid OTP",
		}
		json.NewEncoder(write).Encode(response)
	}
}
func (OtpControll *OtpController) SendOtpHandler(write http.ResponseWriter, request *http.Request) {

	email := request.URL.Query().Get("email")

	result := OtpControll.OtpServe.SendOtp(email)
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(write).Encode(map[string]string{
		"status":  result,
		"message": "OTP is being sent to your email",
	})
}
