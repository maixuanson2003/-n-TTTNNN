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
	Database := AlbumRepo.DB
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
func (AlbumRepo *AlbumRepository) DeleteAlbumById(Id int) error {
	Database := AlbumRepo.DB
	var Album entity.Album
	err := Database.Preload("Song").Preload("Artist").First(&Album, Id).Error

	if err != nil {
		return err
	}
	for _, Song := range Album.Song {
		var SongHandle entity.Song
		errs := Database.
			Preload("SongType").
			Preload("ListenHistory").
			Preload("Review").
			Preload("Artist").
			Preload("User").
			Preload("PlayList").
			Preload("Collection").
			First(&SongHandle, Song.ID).Error
		if errs != nil {
			return errs
		}
		errsDelete := Database.
			Select("SongType", "ListenHistory", "Review", "Artist", "User", "PlayList", "Collection").
			Delete(&SongHandle).Error
		if errsDelete != nil {
			return errsDelete
		}

	}
	errors := Database.Select("Song", "Artist").Delete(&Album).Error
	if errors != nil {
		return errors
	}
	return nil

}
