package otpservice

import (
	"fmt"
	"math/rand"
	"ten_module/internal/Config"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
	"time"
)

type OtpService struct {
	OtpRepo *repository.OtpRepository
}

var OtpServe *OtpService

func InitOtpService() {
	OtpServe = &OtpService{
		OtpRepo: repository.OtpRepo,
	}
}

type OtpServiceInterface interface {
	CheckOtp(otp string) bool
	SendOtp() string
}

func (OtpServe *OtpService) CheckOtp(otp string) bool {
	otpRepo := OtpServe.OtpRepo

	Otp := otpRepo.GetOtp(otp)
	if Otp.Otp == "" || Otp.Otp != otp {
		return false
	}
	return true

}
func GenerateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}
func (OtpServe *OtpService) SendOtp(Email string) string {
	go func() {
		otp := GenerateOTP()
		newOtp := entity.Otp{
			Otp:       otp,
			Create_at: time.Now(),
		}
		OtpServe.OtpRepo.CreateOtp(newOtp)
		Config.SendEmail(Email, "xac thuc otp", otp)
	}()
	return "success"
}
