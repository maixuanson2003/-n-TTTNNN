package repository

import (
	"log"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type PlayListRepository struct {
	DB *gorm.DB
}

var PlayListRepo *PlayListRepository

func InitPlayListRepository() {
	PlayListRepo = &PlayListRepository{
		DB: database.Database,
	}
}

type PlayListRepositoryInterface interface {
	FindAll() ([]entity.PlayList, error)
	GetPlayListById(Id int) (entity.PlayList, error)
	CreatePlayList(Artist entity.PlayList) error
	UpdatePlayList(Artist entity.PlayList, id int) error
	DeletePlayListById(Id int) error
}

func (PlayListRepo *PlayListRepository) FindAll() ([]entity.PlayList, error) {
	Database := ArtistRepo.DB
	var PlayList []entity.PlayList
	err := Database.Model(&entity.Artist{}).Find(&PlayList).Error
	if err != nil {
		return nil, err
	}
	return PlayList, nil
}
func (PlayListRepo *PlayListRepository) GetPlayListById(Id int) (entity.PlayList, error) {
	Database := ArtistRepo.DB
	var PlayList entity.PlayList
	err := Database.Model(&entity.PlayList{}).Preload("Song").Where("id=?", Id).First(&PlayList).Error
	if err != nil {
		return entity.PlayList{}, err
	}
	return PlayList, nil
}
func (PlayListRepo *PlayListRepository) CreatePlayList(PlayList entity.PlayList) error {
	Database := ArtistRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Create(&PlayList).Error
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
func (PlayListRepo *PlayListRepository) UpdatePlayList(PlayList entity.PlayList, id int) error {
	Database := PlayListRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Where("id=?", id).Save(&PlayList).Error
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
