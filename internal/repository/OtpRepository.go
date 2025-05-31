package repository

import (
	"log"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type OtpRepository struct {
	DB *gorm.DB
}

var OtpRepo *OtpRepository

func InitOtpRepository() {
	OtpRepo = &OtpRepository{
		DB: database.Database,
	}
}

type OtpRepositoryInterface interface {
	FindAll() ([]entity.ListenHistory, error)
	GetHistoryById(Id int) (entity.ListenHistory, error)
	CreateHistory(History entity.ListenHistory) error
}

func (OtpRepo *OtpRepository) CreateOtp(otp entity.Otp) error {
	Database := OtpRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Create(&otp).Error
		if err != nil {
			return err

		}
		return nil

	})
	if errs != nil {
		log.Print(errs)
		return errs
	}
	return nil
}
func (OtpRepo *OtpRepository) GetOtp(otp string) entity.Otp {
	Database := OtpRepo.DB
	var Otp entity.Otp
	errs := Database.Where("otp = ?", otp).Find(&Otp).Error
	if errs != nil {
		log.Print(errs)
		return entity.Otp{}
	}
	return Otp

}
