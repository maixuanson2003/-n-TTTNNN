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
	err := Database.Model(&entity.Artist{}).Preload("Song").Preload("Album").Where("id=?", Id).First(&Artist).Error
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
func (ArtistRepo *ArtistRepository) UpdateAritst(Artist entity.Artist, id int) error {
	Database := CollectionRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Where("id=?", id).Save(&Artist).Error
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
func (ArtistRepo *ArtistRepository) SearchArtist(Keyword string) ([]entity.Artist, error) {
	Database := ArtistRepo.DB
	var Artist []entity.Artist
	err := Database.Model(&entity.Artist{}).Where("name LIKE ?", "%"+Keyword+"%").Find(&Artist).Error
	if err != nil {
		return nil, err
	}
	return Artist, nil
}
func (ArtistRepo *ArtistRepository) FilterArtist(CountryId int) ([]entity.Artist, error) {
	Database := ArtistRepo.DB
	var Artist []entity.Artist
	err := Database.Model(&entity.Artist{}).Where("country_id = ?", CountryId).Find(&Artist).Error
	if err != nil {
		return nil, err
	}
	return Artist, nil
}
func (ArtistRepo *ArtistRepository) DeleteArtist(artistId int) error {
	Database := ArtistRepo.DB
	var artist entity.Artist
	err := Database.Preload("Song").Preload("Album").First(&artist, artistId).Error
	if err != nil {
		return err
	}
	errors := Database.Select("Song", "Album").Delete(&artist).Error
	if errors != nil {
		return errors
	}
	return nil
}
