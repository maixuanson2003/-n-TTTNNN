package repository

import (
	"log"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type AlbumRepository struct {
	DB *gorm.DB
}

var AlbumRepo *AlbumRepository

func InitAlbumRepository() {
	AlbumRepo = &AlbumRepository{
		DB: database.Database,
	}
}

type AlbumRepositoryInterface interface {
	FindAll() ([]entity.Album, error)
	GetAlbumById(Id int) (entity.Album, error)
	CreateAlbum(Album entity.Album) error
	UpdateAlbum(Album entity.Album, id int) error
	DeleteAlbumById(Id int) error
}

func (AlbumRepo *AlbumRepository) FindAll() ([]entity.Album, error) {
	Database := ArtistRepo.DB
	var Album []entity.Album
	err := Database.Model(&entity.Album{}).Preload("Song").Preload("Artist").Find(&Album).Error
	if err != nil {
		return nil, err
	}
	return Album, nil

}
func (AlbumRepo *AlbumRepository) GetAlbumById(Id int) (entity.Album, error) {
	Database := ArtistRepo.DB
	var Album entity.Album
	err := Database.Model(&entity.Album{}).Preload("Song").Preload("Artist").Where("id=?", Id).First(&Album).Error
	if err != nil {
		return entity.Album{}, err
	}
	return Album, nil
}
func (AlbumRepo *AlbumRepository) CreateAlbum(Album entity.Album) (int, error) {
	Database := AlbumRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Create(&Album).Error
		if err != nil {
			return err

		}
		return nil

	})
	if errs != nil {
		log.Print(errs)
		return -1, errs
	}
	return Album.ID, nil
}
func (AlbumRepo *AlbumRepository) UpdateAlbum(Album entity.Album, id int) error {
	Database := AlbumRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Where("id=?", id).Save(&Album).Error
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
