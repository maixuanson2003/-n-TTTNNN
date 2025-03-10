package repository

import (
	"log"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type ArtistRepository struct {
	DB *gorm.DB
}

var ArtistRepo *ArtistRepository

func InitArtistRepository() {
	ArtistRepo = &ArtistRepository{
		DB: database.Database,
	}
}

type ArtistRepositoryInterface interface {
	FindAll() ([]entity.Artist, error)
	GetArtistById(Id int) (entity.Artist, error)
	CreateArtist(Artist entity.Artist) error
	UpdateArtist(Artist entity.Artist, id string) error
	DeleteArtistById(Id int) error
}

func (ArtistRepo *ArtistRepository) FindAll() ([]entity.Artist, error) {
	Database := ArtistRepo.DB
	var Artist []entity.Artist
	err := Database.Model(&entity.Artist{}).Find(&Artist).Error
	if err != nil {
		return nil, err
	}
	return Artist, nil
}
func (ArtistRepo *ArtistRepository) GetArtistById(Id int) (entity.Artist, error) {
	Database := ArtistRepo.DB
	var Artist entity.Artist
	err := Database.Model(&entity.Artist{}).Where("id=?", Id).First(&Artist).Error
	if err != nil {
		return entity.Artist{}, err
	}
	return Artist, nil
}
func (ArtistRepo *ArtistRepository) CreateArtist(Artist entity.Artist) error {
	Database := ArtistRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Create(&Artist).Error
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
