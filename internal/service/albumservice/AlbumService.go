package albumservice

import (
	"log"
	"mime/multipart"
	"ten_module/internal/Config"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
	"ten_module/internal/service/songservice"
	"time"
)

type AlbumSerivce struct {
	ArtistRepo   *repository.ArtistRepository
	AlbumRepo    *repository.AlbumRepository
	SongRepo     *repository.SongRepository
	SongTypeRepo *repository.SongTypeRepository
	CountryRepo  *repository.CountryRepository
}
type MessageResponse struct {
	Message string
	Status  string
}

var AlbumServe *AlbumSerivce

func InitAlbumSerivce() {
	AlbumServe = &AlbumSerivce{
		ArtistRepo:   repository.ArtistRepo,
		AlbumRepo:    repository.AlbumRepo,
		SongRepo:     repository.SongRepo,
		SongTypeRepo: repository.SongTypeRepo,
		CountryRepo:  repository.CountryRepo,
	}
}

type AlbumServiceInterface interface {
	GetListAlbum() ([]response.AlbumResponse, error)
	GetAlbumById(Id int) (response.AlbumResponse, error)
	GetAlbumByArtist(artistId int) ([]response.AlbumResponse, error)
	CreateAlbum(AlbumReq request.AlbumRequest, SongFileAlum []request.SongFileAlbum) (MessageResponse, error)
	UpdateAlbum(AlbumReq request.AlbumRequest) (MessageResponse, error)
}

func MapArtistEntityToResponse(Artist entity.Artist, NameCountry string) response.ArtistResponse {
	return response.ArtistResponse{
		ID:          Artist.ID,
		Name:        Artist.Name,
		BirthDay:    Artist.BirthDay,
		Description: Artist.Description,
		Country:     NameCountry,
	}
}

func AlbumEntityMapToAlbumResponse(Album entity.Album, countryRepo *repository.CountryRepository, songrepo *repository.SongRepository) response.AlbumResponse {
	SongEntity := Album.Song
	Artist := Album.Artist
	ArtistResponse := []response.ArtistResponse{}
	SongResponse := []response.SongResponseAlbum{}
	for _, SongItem := range SongEntity {
		songItem, _ := songrepo.GetSongById(SongItem.ID)
		songArtistResponses := []response.ArtistResponse{}
		log.Print(songItem.Artist)
		for _, artist := range songItem.Artist {
			country, err := countryRepo.GetCountryById(artist.CountryId)
			if err != nil {
				log.Printf("Lỗi khi lấy quốc gia của nghệ sĩ ID %d: %v", artist.ID, err)
				continue
			}
			songArtistResponses = append(songArtistResponses, MapArtistEntityToResponse(artist, country.CountryName))
		}

		// Map thể loại của bài hát
		songTypeResponses := []response.SongTypeResponse{}
		for _, songType := range songItem.SongType {
			songTypeResponses = append(songTypeResponses, response.SongTypeResponse{Id: songType.ID, Type: songType.Type})
		}

		SongResponse = append(SongResponse, response.SongResponseAlbum{
			ID:           SongItem.ID,
			NameSong:     SongItem.NameSong,
			Description:  SongItem.Description,
			ReleaseDay:   SongItem.ReleaseDay,
			CreateDay:    SongItem.CreateDay,
			UpdateDay:    SongItem.UpdateDay,
			Point:        SongItem.Point,
			LikeAmount:   SongItem.LikeAmount,
			CountryId:    SongItem.CountryId,
			Status:       SongItem.Status,
			ListenAmout:  SongItem.ListenAmout,
			AlbumId:      SongItem.AlbumId,
			SongResource: SongItem.SongResource,
			Artist:       songArtistResponses,
			SongType:     songTypeResponses,
		})
	}
	for _, ArtistItem := range Artist {
		Country, ErrorToGetCountry := countryRepo.GetCountryById(ArtistItem.CountryId)
		if ErrorToGetCountry != nil {
			log.Print(ErrorToGetCountry)
			return response.AlbumResponse{}
		}
		ArtistResponse = append(ArtistResponse, MapArtistEntityToResponse(ArtistItem, Country.CountryName))
	}
	return response.AlbumResponse{
		ID:          Album.ID,
		NameAlbum:   Album.NameAlbum,
		Description: Album.Description,
		ReleaseDay:  Album.ReleaseDay,
		CreateDay:   Album.CreateDay,
		UpdateDay:   Album.UpdateDay,
		ArtistOwner: Album.ArtistOwner,
		Image:       Album.Image,
		Song:        SongResponse,
		Artist:      ArtistResponse,
	}
}
func (AlbumServe *AlbumSerivce) UpdateSongAlbum(AlbumId *int, songReqs []request.SongRequestUpdate, songFiles []request.SongFile) {
	album, err := AlbumServe.AlbumRepo.GetAlbumById(*AlbumId)
	if err != nil {
		log.Printf("Không tìm thấy album ID %d: %v", *AlbumId, err)
		return
	}
	existingSongs := make(map[int]entity.Song)
	for _, song := range album.Song {
		existingSongs[song.ID] = song
	}
	keepSongIds := make(map[int]bool)
	for index, songReq := range songReqs {
		if songReq.ID != nil && *songReq.ID > 0 {

			if existingSong, exists := existingSongs[*songReq.ID]; exists {
				keepSongIds[*songReq.ID] = true
				updatedSong := existingSong
				updatedSong.NameSong = songReq.NameSong
				updatedSong.Description = songReq.Description
				updatedSong.Point = songReq.Point
				updatedSong.Status = songReq.Status
				updatedSong.CountryId = songReq.CountryId
				updatedSong.ReleaseDay = songReq.ReleaseDay
				updatedSong.UpdateDay = time.Now()
				newSongTypes := []entity.SongType{}
				for _, typeId := range songReq.SongType {
					songType, err := AlbumServe.SongTypeRepo.GetSongTypeById(typeId)
					if err != nil {
						log.Printf("Không tìm thấy thể loại ID %d: %v", typeId, err)
						continue
					}
					newSongTypes = append(newSongTypes, songType)
				}
				updatedSong.SongType = newSongTypes

				newArtists := []entity.Artist{}
				for _, artistId := range songReq.Artist {
					artist, err := AlbumServe.ArtistRepo.GetArtistById(artistId)
					if err != nil {
						log.Printf("Không tìm thấy nghệ sĩ ID %d: %v", artistId, err)
						continue
					}
					newArtists = append(newArtists, artist)
				}
				updatedSong.Artist = newArtists

				if index < len(songFiles) && songFiles[index].File != nil {

					resourceUrl, err := Config.HandleUpLoadFile(songFiles[index].File, songReq.NameSong)
					if err != nil {
						log.Printf("Upload bài hát thất bại: %v", err)
					} else {
						updatedSong.SongResource = resourceUrl
					}
				}
				err := AlbumServe.SongRepo.UpdateSong(updatedSong, *songReq.ID)
				if err != nil {
					log.Printf("Update bài hát ID %d thất bại: %v", *songReq.ID, err)
				} else {
					log.Printf("Update bài hát ID %d thành công", *songReq.ID)
				}
			} else {
				log.Printf("Không tìm thấy bài hát ID %d trong album", *songReq.ID)
			}
		} else {
			songEntity := entity.Song{}
			songEntity.NameSong = songReq.NameSong
			songEntity.Description = songReq.Description
			songEntity.Point = songReq.Point
			songEntity.Status = songReq.Status
			songEntity.CountryId = songReq.CountryId
			songEntity.ReleaseDay = songReq.ReleaseDay
			songEntity.CreateDay = time.Now()
			songEntity.UpdateDay = time.Now()
			songEntity.AlbumId = AlbumId
			newSongTypes := []entity.SongType{}
			for _, typeId := range songReq.SongType {
				songType, err := AlbumServe.SongTypeRepo.GetSongTypeById(typeId)
				if err != nil {
					log.Printf("Không tìm thấy thể loại ID %d: %v", typeId, err)
					continue
				}
				newSongTypes = append(newSongTypes, songType)
			}
			songEntity.SongType = newSongTypes
			newArtists := []entity.Artist{}
			for _, artistId := range songReq.Artist {
				artist, err := AlbumServe.ArtistRepo.GetArtistById(artistId)
				if err != nil {
					log.Printf("Không tìm thấy nghệ sĩ ID %d: %v", artistId, err)
					continue
				}
				newArtists = append(newArtists, artist)
			}
			songEntity.Artist = newArtists
			if index < len(songFiles) && songFiles[index].File != nil {
				resourceUrl, err := Config.HandleUpLoadFile(songFiles[index].File, songReq.NameSong)
				if err != nil {
					log.Printf("Upload bài hát mới thất bại: %v", err)
				} else {
					songEntity.SongResource = resourceUrl
				}
			}
			err := AlbumServe.SongRepo.CreateSong(songEntity)
			if err != nil {
				log.Printf("Tạo bài hát mới thất bại: %v", err)
			} else {
				log.Printf("Tạo bài hát mới '%s' thành công", songReq.NameSong)
			}
		}
	}
	deleteIds := []int{}
	for _, song := range album.Song {
		if !keepSongIds[song.ID] {
			deleteIds = append(deleteIds, song.ID)
		}
	}

	for _, id := range deleteIds {
		log.Printf("Xóa bài hát ID %d", id)
		err := AlbumServe.SongRepo.DeleteSongById(id)
		if err != nil {
			log.Printf("Xóa bài hát ID %d thất bại: %v", id, err)
		} else {
			log.Printf("Xóa bài hát ID %d thành công", id)
		}
	}

	log.Printf("Hoàn thành cập nhật album ID %d", *AlbumId)
}
func (AlbumServe *AlbumSerivce) CreateAlbum(AlbumReq request.AlbumRequest, SongFileAlum []request.SongFileAlbum, File *multipart.FileHeader) (MessageResponse, error) {
	AlbumRepo := AlbumServe.AlbumRepo
	SongResourceHasmap := map[string]multipart.File{}
	for _, SongValue := range SongFileAlum {
		SongResourceHasmap[SongValue.SongName] = SongValue.File
	}
	NewAlbum := entity.Album{
		NameAlbum:   AlbumReq.NameAlbum,
		Description: AlbumReq.Description,
		ReleaseDay:  AlbumReq.ReleaseDay,
		ArtistOwner: AlbumReq.ArtistOwner,
		CreateDay:   time.Now(),
		UpdateDay:   time.Now(),
	}
	AlbumId, ErrorToCreateAlbum := AlbumRepo.CreateAlbum(NewAlbum)
	if ErrorToCreateAlbum != nil {
		log.Print(ErrorToCreateAlbum)
		return MessageResponse{
			Message: "failed to create album",
			Status:  "Failed",
		}, ErrorToCreateAlbum
	}

	go func() {
		err := AlbumServe.processAlbumBackground(AlbumId, AlbumReq, SongFileAlum, File)
		if err != nil {
			log.Print("loi gui file cho cloudinary")
		}
	}()

	return MessageResponse{
		Message: "success to create album",
		Status:  "Success",
	}, nil

}
func (AlbumServe *AlbumSerivce) processAlbumBackground(AlbumId int, AlbumReq request.AlbumRequest, SongFileAlum []request.SongFileAlbum, File *multipart.FileHeader) error {
	AlbumRepo := AlbumServe.AlbumRepo
	ArtistRepo := AlbumServe.ArtistRepo
	SongTypeRepo := AlbumServe.SongTypeRepo

	SongResourceMap := map[string]multipart.File{}
	for _, song := range SongFileAlum {
		SongResourceMap[song.SongName] = song.File
	}

	AlbumEntiy, err := AlbumRepo.GetAlbumById(AlbumId)
	if err != nil {
		return err
	}

	for _, SongReq := range AlbumReq.Song {
		var SongTypeArray []entity.SongType
		for _, id := range SongReq.SongType {
			entity, err := SongTypeRepo.GetSongTypeById(id)
			if err != nil {
				return err
			}
			SongTypeArray = append(SongTypeArray, entity)
		}
		var ArtistArray []entity.Artist
		for _, id := range SongReq.Artist {
			entity, err := ArtistRepo.GetArtistById(id)
			if err != nil {
				return err
			}
			ArtistArray = append(ArtistArray, entity)
		}
		file := SongResourceMap[SongReq.NameSong]
		songResource, err := Config.HandleUpLoadFile(file, SongReq.NameSong)

		if err != nil {
			return err
		}
		AlbumEntiy.Song = append(AlbumEntiy.Song,
			songservice.SongReqMapToSongEntity(SongReq, songResource, SongTypeArray, ArtistArray),
		)
	}
	imageResource, err := Config.HandleUploadImage(File, AlbumEntiy.NameAlbum)

	for _, id := range AlbumReq.Artist {
		entity, err := ArtistRepo.GetArtistById(id)
		if err != nil {
			return err
		}
		AlbumEntiy.Artist = append(AlbumEntiy.Artist, entity)
	}
	AlbumEntiy.Image = imageResource

	err = AlbumRepo.UpdateAlbum(AlbumEntiy, AlbumId)
	return err
}

func (AlbumServe *AlbumSerivce) GetListAlbum() ([]response.AlbumResponse, error) {
	AlbumRepo := AlbumServe.AlbumRepo
	AlbumList, ErrorToGetListAlbum := AlbumRepo.FindAll()
	if ErrorToGetListAlbum != nil {
		log.Print(ErrorToGetListAlbum)
		return nil, ErrorToGetListAlbum

	}
	AlbumListResponse := []response.AlbumResponse{}
	for _, AlbumItem := range AlbumList {
		AlbumListResponse = append(AlbumListResponse, AlbumEntityMapToAlbumResponse(AlbumItem, AlbumServe.CountryRepo, AlbumServe.SongRepo))
	}
	return AlbumListResponse, nil
}
func (AlbumServe *AlbumSerivce) GetAlbumById(Id int) (response.AlbumResponse, error) {
	AlbumRepo := AlbumServe.AlbumRepo
	AlbumItem, ErrorToGetAlbum := AlbumRepo.GetAlbumById(Id)
	if ErrorToGetAlbum != nil {
		log.Print(ErrorToGetAlbum)
		return response.AlbumResponse{}, ErrorToGetAlbum

	}
	AlbumRespone := AlbumEntityMapToAlbumResponse(AlbumItem, AlbumServe.CountryRepo, AlbumServe.SongRepo)
	return AlbumRespone, nil
}
func (AlbumServe *AlbumSerivce) GetAlbumByArtist(artistId int) ([]response.AlbumResponse, error) {
	ArtistRepo := AlbumServe.ArtistRepo
	AlbumRepo := AlbumServe.AlbumRepo
	ArtistItem, ErrorToGetArtist := ArtistRepo.GetArtistById(artistId)
	if ErrorToGetArtist != nil {
		log.Print(ErrorToGetArtist)
		return nil, ErrorToGetArtist
	}
	AlbumList := ArtistItem.Album
	AlbumListResponse := []response.AlbumResponse{}
	for _, AlbumItem := range AlbumList {
		Album, Error := AlbumRepo.GetAlbumById(AlbumItem.ID)
		if Error != nil {
			log.Print(Error)
			return nil, Error
		}
		AlbumListResponse = append(AlbumListResponse, AlbumEntityMapToAlbumResponse(Album, AlbumServe.CountryRepo, AlbumServe.SongRepo))
	}
	return AlbumListResponse, nil
}
func (AlbumServe *AlbumSerivce) DeleteAlbum(albumId int) (MessageResponse, error) {
	AlbumRepo := AlbumServe.AlbumRepo
	err := AlbumRepo.DeleteAlbumById(albumId)
	if err != nil {
		return MessageResponse{}, err
	}
	return MessageResponse{
		Message: "success",
		Status:  "Success",
	}, nil
}
func (AlbumServe *AlbumSerivce) processAlbumImage(Album entity.Album, Albumrepo *repository.AlbumRepository, File *multipart.FileHeader) error {
	imageResource, err := Config.HandleUploadImage(File, Album.NameAlbum)
	if err != nil {
		return err
	}
	Album.Image = imageResource
	errs := Albumrepo.UpdateAlbum(Album, Album.ID)
	return errs

}
func (AlbumServe *AlbumSerivce) UpdateAlbum(AlbumReq request.AlbumUpdate, Id int, File *multipart.FileHeader) (MessageResponse, error) {
	AlbumRepo := AlbumServe.AlbumRepo

	AlbumItem, err := AlbumRepo.GetAlbumById(Id)
	if err != nil {
		log.Printf("Không tìm thấy album ID %d: %v", Id, err)
		return MessageResponse{}, err
	}

	AlbumItem.NameAlbum = AlbumReq.NameAlbum
	AlbumItem.Description = AlbumReq.Description
	AlbumItem.ReleaseDay = AlbumReq.ReleaseDay
	AlbumItem.ArtistOwner = AlbumReq.ArtistOwner

	newArtists := []entity.Artist{}
	for _, artistID := range AlbumReq.Artist {
		artist, e := AlbumServe.ArtistRepo.GetArtistById(artistID)
		if e != nil {
			log.Printf("Không tìm thấy nghệ sĩ ID %d: %v", artistID, e)
			continue
		}
		newArtists = append(newArtists, artist)
	}
	AlbumItem.Artist = newArtists

	/* === 4. Lưu album ======================================== */
	if err := AlbumRepo.UpdateAlbum(AlbumItem, Id); err != nil {
		log.Print(err)
		return MessageResponse{
			Message: "failed",
			Status:  "Failed",
		}, err
	}

	go func() {
		if File == nil {
			return
		}
		if err := AlbumServe.processAlbumImage(AlbumItem, AlbumRepo, File); err != nil {
			log.Print(err)
		}
	}()

	return MessageResponse{
		Message: "success",
		Status:  "Success",
	}, nil

}
