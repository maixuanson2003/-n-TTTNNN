package repository

import (
	"log"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type ListenHistoryRepo struct {
	DB *gorm.DB
}

var ListenRepo *ListenHistoryRepo

func InitListenHistoryRepo() {
	ListenRepo = &ListenHistoryRepo{
		DB: database.Database,
	}
}

type ListenHistoryRepoInterface interface {
	FindAll() ([]entity.ListenHistory, error)
	GetHistoryById(Id int) (entity.ListenHistory, error)
	CreateHistory(History entity.ListenHistory) error
}

func (HisRepo *ListenHistoryRepo) CreateHistory(History entity.ListenHistory) error {
	Database := HisRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Create(&History).Error
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
func (HisRepo *ListenHistoryRepo) CountNumberSongId(songId int) (int64, error) {
	Database := HisRepo.DB
	var count int64
	err := Database.Model(&entity.ListenHistory{}).
		Where("song_id = ?", songId).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}
