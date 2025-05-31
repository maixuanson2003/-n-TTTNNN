package otpservice

import (
	"ten_module/internal/repository"
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
}

func (OtpServe *OtpService) CheckOtp(otp string) bool {
	otpRepo := OtpServe.OtpRepo

	Otp := otpRepo.GetOtp(otp)
	if Otp.Otp == "" || Otp.Otp != otp {
		return false
	}
	return true

}
