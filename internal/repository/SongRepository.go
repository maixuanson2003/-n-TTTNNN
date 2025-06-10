package repository

import (
	"log"
	"strings"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"
	gemini "ten_module/internal/Helper/openAi"

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

	err := Database.Model(&entity.Song{}).Preload("SongType").Preload("ListenHistory").Preload("Artist").Find(&Song).Error
	if err != nil {
		return nil, err
	}
	return Song, nil

}
func (songRepository *SongRepository) Paginate(page int) ([]entity.Song, error) {
	Database := songRepository.DB
	var Song []entity.Song
	batchSize := 6
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
func (songRepository *SongRepository) DeleteSongsByAlbumId(albumId int) error {
	db := songRepository.DB

	return db.Transaction(func(tx *gorm.DB) error {
		var songs []entity.Song
		if err := tx.Where("album_id = ?", albumId).Find(&songs).Error; err != nil {
			return err
		}

		for _, song := range songs {
			songID := song.ID

			if err := tx.Where("song_id = ?", songID).Delete(&entity.ListenHistory{}).Error; err != nil {
				return err
			}
			if err := tx.Where("song_id = ?", songID).Delete(&entity.Review{}).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM song_song_types WHERE song_id = ?", songID).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM song_artists WHERE song_id = ?", songID).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM user_likes WHERE song_id = ?", songID).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM play_list_songs WHERE song_id = ?", songID).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM collection_songs WHERE song_id = ?", songID).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("album_id = ?", albumId).Delete(&entity.Song{}).Error; err != nil {
			return err
		}

		return nil
	})
}
func (r *SongRepository) SearchOrRecommendSongs(query *gemini.MusicQuery) ([]entity.Song, error) {
	var songs []entity.Song

	db := r.DB.Model(&entity.Song{}).
		Preload("Artist").
		Preload("SongType").
		Preload("Album")

	if query.Song != "" {
		db = db.Where("LOWER(name_song) LIKE ?", "%"+strings.ToLower(query.Song)+"%")
	}

	if query.Artist != "" {
		db = db.Joins("JOIN song_artists ON song_artists.song_id = songs.id").
			Joins("JOIN artists ON artists.id = song_artists.artist_id").
			Where("LOWER(artists.name) LIKE ?", "%"+strings.ToLower(query.Artist)+"%")
	}

	if query.Album != "" {
		db = db.Joins("JOIN albums ON albums.id = songs.album_id").
			Where("LOWER(albums.name) LIKE ?", "%"+strings.ToLower(query.Album)+"%")
	}

	if query.Genre != "" {
		db = db.Joins("JOIN song_song_types ON song_song_types.song_id = songs.id").
			Joins("JOIN song_types ON song_types.id = song_song_types.song_type_id").
			Where("LOWER(song_types.type) LIKE ?", "%"+strings.ToLower(query.Genre)+"%")
	}

	if query.Keywords != "" {
		db = db.Where("LOWER(songs.name_song) LIKE ? OR LOWER(songs.description) LIKE ?",
			"%"+strings.ToLower(query.Keywords)+"%",
			"%"+strings.ToLower(query.Keywords)+"%")
	}

	switch query.TimeRange {
	case "today":
		db = db.Where("DATE(songs.create_day) = CURDATE()")
	case "week":
		db = db.Where("YEARWEEK(songs.create_day, 1) = YEARWEEK(CURDATE(), 1)")
	case "month":
		db = db.Where("MONTH(songs.create_day) = MONTH(CURDATE()) AND YEAR(create_day) = YEAR(CURDATE())")
	case "year":
		db = db.Where("YEAR(songs.create_day) = YEAR(CURDATE())")
	}

	switch query.SortBy {
	case "most_played":
		db = db.Order("songs.listen_amout DESC")
	case "latest":
		db = db.Order("songs.create_day DESC")
	case "top":
		db = db.Order("songs.point DESC, songs.like_amount DESC")
	case "popular":
		db = db.Order("songs.like_amount DESC")
	default:
		db = db.Order("songs.create_day DESC")
	}
	if query.Intent == "play" {
		db = db.Limit(1)
	} else {
		db = db.Limit(10)
	}
	if err := db.Distinct("songs.*").Find(&songs).Error; err != nil {
		return nil, err
	}

	return songs, nil
}
