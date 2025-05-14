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
func (songRepository *SongRepository) FilterSong(ArtistId []int, TypeId []int) ([]entity.Song, error) {
	Database := songRepository.DB
	var Song []entity.Song
	print(len((ArtistId)))
	for _, id := range ArtistId {
		print(id)
	}
	Query := Database.Model(&entity.Song{}).
		Joins("JOIN song_artists ON song_artists.song_id=songs.id").
		Joins("JOIN artists ON artists.id=song_artists.artist_id").
		Joins("JOIN song_song_types ON song_song_types.song_id=songs.id").
		Joins("JOIN song_types ON song_types.id=song_song_types.song_type_id")

	if len(ArtistId) > 0 {
		Query = Query.Where("song_artists.artist_id IN ?", ArtistId)
	}
	if len(TypeId) > 0 {
		Query = Query.Where("song_song_types.song_type_id IN ?", TypeId)
	}
	err := Query.Preload("Review").Preload("SongType").Preload("Artist").Group("songs.id").Find(&Song).Error
	if err != nil {
		return nil, err
	}
	log.Print(Song)

	return Song, nil
}
func (songRepository *SongRepository) DeleteSongById(id int) error {
	db := songRepository.DB

	return db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("song_id = ?", id).Delete(&entity.ListenHistory{}).Error; err != nil {
			return err
		}

		if err := tx.Where("song_id = ?", id).Delete(&entity.Review{}).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM song_song_types WHERE song_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM song_artists WHERE song_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM user_likes WHERE song_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM play_list_songs WHERE song_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM collection_songs WHERE song_id = ?", id).Error; err != nil {
			return err
		}

		// Cuối cùng xóa bài hát
		if err := tx.Delete(&entity.Song{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}
