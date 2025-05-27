package repository

import (
	"log"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type CollectionRepostiory struct {
	DB *gorm.DB
}
type CollectionRepostioryInterface interface {
	FindAll() ([]entity.Collection, error)
	GetCollectById(Id int) (entity.Collection, error)
	CreateCollect(Collect entity.Collection) error
	UpdateCollect(Collect entity.Collection, id int) error
	DeleteCollectById(Id int) error
	DeleteSong(SongId int, CollectionId int) error
}

var CollectionRepo *CollectionRepostiory

func InitCollectionRepostiory() {
	CollectionRepo = &CollectionRepostiory{
		DB: database.Database,
	}
}
func (CollectionRepo *CollectionRepostiory) FindAll() ([]entity.Collection, error) {
	Database := ArtistRepo.DB
	var Collection []entity.Collection
	err := Database.Model(&entity.Collection{}).Find(&Collection).Error
	if err != nil {
		return nil, err
	}
	return Collection, nil
}
func (CollectionRepo *CollectionRepostiory) GetCollectById(Id int) (entity.Collection, error) {
	Database := CollectionRepo.DB
	var Collection entity.Collection
	err := Database.Model(&entity.Collection{}).Preload("Song").Where("id=?", Id).First(&Collection).Error
	if err != nil {
		return entity.Collection{}, err
	}
	return Collection, nil
}
func (CollectionRepo *CollectionRepostiory) CreateCollect(Collect entity.Collection) error {
	Database := CollectionRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Create(&Collect).Error
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
func (CollectionRepo *CollectionRepostiory) UpdateCollect(Collect entity.Collection, id int) error {
	Database := CollectionRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Where("id=?", id).Save(&Collect).Error
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
func (CollectionRepo *CollectionRepostiory) DeleteSong(SongId int, CollectionId int) error {
	Database := CollectionRepo.DB
	var Collection entity.Collection
	errorToFindCollection := Database.Preload("Song").Where("id=?", CollectionId).First(&Collection).Error
	if errorToFindCollection != nil {
		return errorToFindCollection
	}
	errorToDeleteSongFromCollection := Database.Model(&Collection).Association("Song").Delete(entity.Song{ID: SongId})
	if errorToDeleteSongFromCollection != nil {
		return errorToDeleteSongFromCollection
	}
	return nil
}

func (CollectionRepo *CollectionRepostiory) DeleteCollectById(Id int) error {
	Database := CollectionRepo.DB
	var collection entity.Collection
	err := Database.Preload("Song").First(&collection, Id).Error
	if err != nil {
		return err
	}
	errors := Database.Select("Song").Delete(&collection).Error
	if errors != nil {
		return errors
	}
	return nil
}
