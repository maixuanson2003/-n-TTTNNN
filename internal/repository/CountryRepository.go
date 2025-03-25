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
	err := Database.Model(&entity.Artist{}).Find(&Country).Error
	if err != nil {
		return nil, err
	}
	return Country, nil
}
func (CountryRepo *CountryRepository) GetCountryById(Id int) (entity.Country, error) {
	Database := CountryRepo.DB
	var Country entity.Country
	err := Database.Model(&entity.Album{}).Preload("Song").Preload("Artist").Where("id=?", Id).First(&Country).Error
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
	Database := CollectionRepo.DB
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
