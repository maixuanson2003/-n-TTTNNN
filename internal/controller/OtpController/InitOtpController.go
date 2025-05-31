package otpcontroller

import (
	"ten_module/internal/service/otpservice"
)

func InitOtpControll() {
	otpservice.InitOtpService()
	InitOtpController()
}
