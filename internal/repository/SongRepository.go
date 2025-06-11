package repository

import (
	"log"
	"strings"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"
	"time"

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

func parseCSV(input string) []string {
	var result []string
	for _, item := range strings.Split(input, ",") {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// RecommendSongs xử lý genre, artist, keywords, time_range, sort_by
func (repo *SongRepository) RecommendSongs(
	genre string,
	artist string,
	keywords string,
	timeRange string,
	sortBy string,
) ([]entity.Song, error) {
	var songs []entity.Song

	db := repo.DB.Model(&entity.Song{}).
		Select("songs.*, COUNT(listen_histories.id) as listen_count").
		Joins("LEFT JOIN listen_histories ON songs.id = listen_histories.song_id").
		Group("songs.id").
		Preload("SongType").
		Preload("Artist").
		Preload("Album")

	// Filter by Genre (many-to-many)
	if genre != "" {
		genres := parseCSV(genre)
		log.Print(genres)
		if len(genres) > 0 {
			db = db.
				Joins("JOIN song_song_types ON song_song_types.song_id = songs.id").
				Joins("JOIN song_types ON song_types.id = song_song_types.song_type_id").
				Where("song_types.type IN ?", genres)
		}
	}

	log.Print(artist)
	if artist != "" {
		artists := parseCSV(artist)
		if len(artists) > 0 {
			db = db.
				Joins("JOIN song_artists ON song_artists.song_id = songs.id").
				Joins("JOIN artists ON artists.id = song_artists.artist_id").
				Where("artists.name IN ?", artists)
		}
	}

	// Filter by Keywords (title contains)
	// if keywords != "" {
	// 	db = db.Where("songs.name_song LIKE ?", "%"+keywords+"%")
	// }

	if timeRange != "" {
		now := time.Now()
		var fromTime time.Time

		switch strings.ToLower(timeRange) {
		case "today":
			fromTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		case "week":
			fromTime = now.AddDate(0, 0, -7)
		case "month":
			fromTime = now.AddDate(0, -1, 0)
		case "year":
			fromTime = now.AddDate(-1, 0, 0)
		}

		db = db.Where("listen_histories.listen_day >= ?", fromTime)
	}

	// Sort by
	switch strings.ToLower(sortBy) {
	case "latest":
		db = db.Order("songs.create_day DESC")
	case "most_played":
		db = db.Order("songs.listen_amout DESC")
	case "popular":
		db = db.Order("songs.like_amount DESC")
	}

	// Execute query
	err := db.Find(&songs).Error
	if err != nil {
		return nil, err
	}

	return songs, nil
}
