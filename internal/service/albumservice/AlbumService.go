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

func AlbumEntityMapToAlbumResponse(Album entity.Album, countryRepo *repository.CountryRepository) response.AlbumResponse {
	SongEntity := Album.Song
	Artist := Album.Artist
	ArtistResponse := []response.ArtistResponse{}
	SongResponse := []response.SongResponse{}
	for _, SongItem := range SongEntity {
		SongResponse = append(SongResponse, songservice.SongEntityMapToSongResponse(SongItem))
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
		Song:        SongResponse,
		Artist:      ArtistResponse,
	}
}
func (AlbumServe *AlbumSerivce) CreateAlbum(AlbumReq request.AlbumRequest, SongFileAlum []request.SongFileAlbum) (MessageResponse, error) {
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
		err := AlbumServe.processAlbumBackground(AlbumId, AlbumReq, SongFileAlum)
		if err != nil {
			log.Print("loi gui file cho cloudinary")
		}
	}()

	return MessageResponse{
		Message: "success to create album",
		Status:  "Success",
	}, nil

}
func (AlbumServe *AlbumSerivce) processAlbumBackground(AlbumId int, AlbumReq request.AlbumRequest, SongFileAlum []request.SongFileAlbum) error {
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

	for _, id := range AlbumReq.Artist {
		entity, err := ArtistRepo.GetArtistById(id)
		if err != nil {
			return err
		}
		AlbumEntiy.Artist = append(AlbumEntiy.Artist, entity)
	}

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
		AlbumListResponse = append(AlbumListResponse, AlbumEntityMapToAlbumResponse(AlbumItem, AlbumServe.CountryRepo))
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
	AlbumRespone := AlbumEntityMapToAlbumResponse(AlbumItem, AlbumServe.CountryRepo)
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
		AlbumListResponse = append(AlbumListResponse, AlbumEntityMapToAlbumResponse(Album, AlbumServe.CountryRepo))
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
func (AlbumServe *AlbumSerivce) UpdateAlbum(AlbumReq request.AlbumUpdate, Id int) (MessageResponse, error) {
	AlbumRepo := AlbumServe.AlbumRepo
	AlbumItem, errors := AlbumRepo.GetAlbumById(Id)
	if errors != nil {
		log.Print(errors)
		return MessageResponse{}, errors
	}
	AlbumItem.ArtistOwner = AlbumReq.ArtistOwner
	AlbumItem.Description = AlbumReq.Description
	AlbumItem.NameAlbum = AlbumReq.NameAlbum
	AlbumItem.ReleaseDay = AlbumReq.ReleaseDay
	errorsToUpdate := AlbumRepo.UpdateAlbum(AlbumItem, Id)
	if errorsToUpdate != nil {
		log.Print(errorsToUpdate)
		return MessageResponse{
			Message: "failed",
			Status:  "Failed",
		}, errorsToUpdate
	}
	return MessageResponse{
		Message: "success",
		Status:  "Success",
	}, nil

}
