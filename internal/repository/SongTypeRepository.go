package repository

import (
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type SongTypeRepository struct {
	DB *gorm.DB
}
type SongTypeRepositoryInterface interface {
	FindAll() ([]entity.SongType, error)
	GetSongTypeById(Id int) (entity.SongType, error)
	CreateSongType(Song entity.SongType) error
	UpdateSongType(Song entity.SongType, id string) error
	DeleteSongTypeById(Id int) error
	DeleteAll(User []entity.SongType) error
}

var SongTypeRepo *SongTypeRepository

func InitSongTypeRepository() {
	SongTypeRepo = &SongTypeRepository{
		DB: database.Database,
	}
}
func (songTypeRepo *SongTypeRepository) FindAll() ([]entity.SongType, error) {
	Database := songTypeRepo.DB
	var SongType []entity.SongType
	err := Database.Model(&entity.SongType{}).Find(&SongType).Error
	if err != nil {
		return nil, err
	}
	return SongType, nil
}
func (songTypeRepo *SongTypeRepository) GetSongTypeById(Id int) (entity.SongType, error) {
	Database := songTypeRepo.DB
	var SongType entity.SongType
	err := Database.Model(&entity.SongType{}).Preload("Song").Where("id=?", Id).First(&SongType).Error
	if err != nil {
		return entity.SongType{}, err
	}
	return SongType, nil
}
