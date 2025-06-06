package repository

import (
	"log"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type CountryRepository struct {
	DB *gorm.DB
}

var CountryRepo *CountryRepository

func InitCountryRepository() {
	CountryRepo = &CountryRepository{
		DB: database.Database,
	}
}

type CountryRepositoryInterface interface {
	FindAll() ([]entity.Country, error)
	GetCountryById(Id int) (entity.Country, error)
	CreateCountry(Country entity.Country) error
	UpdateCountry(Country entity.Country, id int) error
	DeleteAlbumById(Id int) error
}

func (CountryRepo *CountryRepository) FindAll() ([]entity.Country, error) {
	Database := CountryRepo.DB
	var Country []entity.Country
	err := Database.Model(&entity.Country{}).Find(&Country).Error
	if err != nil {
		return nil, err
	}
	return Country, nil
}
func (CountryRepo *CountryRepository) GetCountryById(Id int) (entity.Country, error) {
	Database := CountryRepo.DB
	var Country entity.Country
	err := Database.Model(&entity.Country{}).Preload("Song").Preload("Artist").Where("id=?", Id).First(&Country).Error
	if err != nil {
		return entity.Country{}, err
	}
	return Country, nil
}
func (CountryRepo *CountryRepository) CreateCountry(Country entity.Country) error {
	Database := CountryRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Create(&Country).Error
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
func (CountryRepo *CountryRepository) UpdateCountry(Country entity.Country, id int) error {
	Database := CountryRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Where("id=?", id).Save(&Country).Error
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
func (CountryRepo *CountryRepository) DeleteCountryByID(id int) error {
	Database := CountryRepo.DB
	var country entity.Country
	if err := Database.Preload("Song").Preload("Artist").Where("id =?", id).First(&country, id).Error; err != nil {
		return err
	}
	song := country.Song
	artist := country.Artist
	for _, songItem := range song {
		var SongHandle entity.Song
		errs := Database.
			Preload("SongType").
			Preload("ListenHistory").
			Preload("Review").
			Preload("Artist").
			Preload("User").
			Preload("PlayList").
			Preload("Collection").
			First(&SongHandle, songItem.ID).Error
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
	for _, artistItem := range artist {
		var artist entity.Artist
		err := Database.Preload("Song").Preload("Album").First(&artist, artistItem.ID).Error
		if err != nil {
			return err
		}
		errors := Database.Select("Song", "Album").Delete(&artist).Error
		if errors != nil {
			return errors
		}
	}

	return Database.Select("Song", "Artist").Delete(&country).Error
}
