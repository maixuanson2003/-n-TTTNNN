package repository

import (
	"log"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type SongRepository struct {
	DB *gorm.DB
}

var SongRepo *SongRepository

func InitSongRepo() {
	SongRepo = &SongRepository{
		DB: database.Database,
	}
}

type SongRepoInterface interface {
	FindAll() ([]entity.Song, error)
	Paginate(offset int) ([]entity.Song, error)
	GetSongById(Id int) (entity.Song, error)
	CreateSong(Song entity.Song) error
	UpdateSong(Song entity.Song, id string) error
	DeleteSongById(Id int) error
	DeleteAll(User []entity.Song) error
	SearchSongByKey(Keyword string) ([]entity.Song, error)
}

func (songRepository *SongRepository) FindAll() ([]entity.Song, error) {
	Database := songRepository.DB
	var Song []entity.Song

	err := Database.Model(&entity.Song{}).Preload("SongType").Preload("Artist").Find(&Song).Error
	if err != nil {
		return nil, err
	}
	return Song, nil

}
func (songRepository *SongRepository) Paginate(page int) ([]entity.Song, error) {
	Database := songRepository.DB
	var Song []entity.Song
	batchSize := 10
	offset := (page - 1) * batchSize
	err := Database.Model(&entity.Song{}).Limit(batchSize).Offset(offset).Preload("SongType").Preload("Artist").Find(&Song).Error
	if err != nil {
		return nil, err
	}
	return Song, nil
}
func (songRepository *SongRepository) GetSongById(Id int) (entity.Song, error) {
	Database := songRepository.DB
	var Song entity.Song
	err := Database.Model(&entity.Song{}).Preload("Review").Preload("SongType").Preload("Artist").Where("id=?", Id).First(&Song).Error
	if err != nil {
		return entity.Song{}, err
	}
	return Song, nil

}
func (songRepository *SongRepository) CreateSong(Song entity.Song) error {
	Database := songRepository.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Create(&Song).Error
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
func (songRepository *SongRepository) UpdateSong(Song entity.Song, id int) error {
	Database := songRepository.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Where("id=?", id).Save(&Song).Error
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
func (songRepository *SongRepository) SearchSongByKey(Keyword string) ([]entity.Song, error) {
	Database := songRepository.DB
	var Song []entity.Song
	err := Database.Model(&entity.Song{}).Preload("Review").Preload("SongType").Preload("Artist").Where("name_song LIKE ?", "%"+Keyword+"%").Find(&Song).Error
	if err != nil {
		return nil, err
	}
	return Song, nil
}
