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
	db := songRepository.DB.Session(&gorm.Session{NewDB: true}).Debug()

	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Lấy thông tin (để xoá file sau commit nếu cần)
		var song entity.Song
		if err := tx.First(&song, id).Error; err != nil {
			return err
		}

		// 2. Xoá các bảng pivot many‑to‑many
		pivotSQL := []string{
			"DELETE FROM song_song_types   WHERE song_id = ?",
			"DELETE FROM song_artists      WHERE song_id = ?",
			"DELETE FROM user_likes        WHERE song_id = ?",
			"DELETE FROM play_list_songs   WHERE song_id = ?",
			"DELETE FROM collection_songs  WHERE song_id = ?",
		}
		for _, q := range pivotSQL {
			if err := tx.Exec(q, id).Error; err != nil {
				return err
			}
		}

		// 3. Xoá bảng con has‑many
		if err := tx.Where("song_id = ?", id).Delete(&entity.ListenHistory{}).Error; err != nil {
			return err
		}
		if err := tx.Where("song_id = ?", id).Delete(&entity.Review{}).Error; err != nil {
			return err
		}

		// 4. Xoá bài hát (hard delete)
		if err := tx.Unscoped().Delete(&entity.Song{}, id).Error; err != nil {
			return err
		}

		log.Printf("✅ Đã xoá bài hát ID %d", id)
		// (tuỳ chọn) Sau commit xoá file vật lý
		// _ = Config.DeleteFile(song.SongResource)
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
	genre, artist, country, keywords, timeRange, sortBy string,
) ([]entity.Song, error) {

	var songs []entity.Song

	db := repo.DB.Model(&entity.Song{}).
		Select("songs.*, COUNT(listen_histories.id) AS listen_count").
		Joins("LEFT JOIN listen_histories ON songs.id = listen_histories.song_id").
		Group("songs.id").
		Preload("SongType").
		Preload("Artist").
		Preload("Album").
		Preload("Country")

	if genre != "" {
		if g := parseCSV(genre); len(g) > 0 {
			db = db.
				Joins("JOIN song_song_types ON song_song_types.song_id = songs.id").
				Joins("JOIN song_types ON song_types.id = song_song_types.song_type_id").
				Where("song_types.type IN ?", g)
		}
	}
	if artist != "" {
		if a := parseCSV(artist); len(a) > 0 {
			db = db.
				Joins("JOIN song_artists ON song_artists.song_id = songs.id").
				Joins("JOIN artists ON artists.id = song_artists.artist_id").
				Where("artists.name IN ?", a)
		}
	}
	if country != "" {
		if c := parseCSV(country); len(c) > 0 {
			// tuỳ cấu trúc DB: nếu songs có country_id
			db = db.
				Joins("JOIN countries ON countries.id = songs.country_id").
				Where("countries.country_name IN ?", c)
		}
	}

	if timeRange != "" {
		now := time.Now()
		var from time.Time
		switch strings.ToLower(timeRange) {
		case "today":
			from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		case "week":
			from = now.AddDate(0, 0, -7)
		case "month":
			from = now.AddDate(0, -1, 0)
		case "year":
			from = now.AddDate(-1, 0, 0)
		}
		db = db.Where("listen_histories.listen_day >= ?", from)
	}

	/* ---------- Sort ---------- */
	switch strings.ToLower(sortBy) {
	case "latest":
		db = db.Order("songs.create_day DESC")
	case "most_played":
		db = db.Order("songs.listen_amout DESC")
	case "popular":
		db = db.Order("songs.like_amount DESC")
	case "top":
		db = db.Order("listen_count DESC")
	}
	if err := db.Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}
