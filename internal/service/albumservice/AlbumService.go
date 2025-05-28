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
	ArtistRepo := AlbumServe.ArtistRepo
	SongTypeRepo := AlbumServe.SongTypeRepo
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
	AlbumEntiy, ErrorToGetAlbum := AlbumRepo.GetAlbumById(AlbumId)
	if ErrorToGetAlbum != nil {
		log.Print(ErrorToGetAlbum)
		return MessageResponse{}, ErrorToGetAlbum

	}

	SongRequestArray := AlbumReq.Song
	for _, SongRequestItem := range SongRequestArray {
		SongTypeId := SongRequestItem.SongType
		SongTypeArray := []entity.SongType{}
		ArtistId := SongRequestItem.Artist
		ArtistArray := []entity.Artist{}
		for _, SongType := range SongTypeId {
			SongTypeEntity, ErrorToGetSongType := SongTypeRepo.GetSongTypeById(SongType)
			if ErrorToGetSongType != nil {
				log.Print(ErrorToGetSongType)
				return MessageResponse{
					Message: "failed to get song type",
					Status:  "Failed",
				}, ErrorToGetSongType
			}
			SongTypeArray = append(SongTypeArray, SongTypeEntity)

		}
		for _, Artist := range ArtistId {
			AritstEntity, ErrorToGetArtist := ArtistRepo.GetArtistById(Artist)
			if ErrorToGetArtist != nil {
				log.Print(ErrorToGetArtist)
				return MessageResponse{
					Message: "failed to get artist for song",
					Status:  "Failed",
				}, ErrorToGetArtist
			}
			ArtistArray = append(ArtistArray, AritstEntity)
		}
		SongResource, ErrorToUploadFile := Config.HandleUpLoadFile(SongResourceHasmap[SongRequestItem.NameSong], SongRequestItem.NameSong)
		if ErrorToUploadFile != nil {
			log.Print(ErrorToUploadFile)
			return MessageResponse{
				Message: "failed to create resource",
				Status:  "Failed",
			}, ErrorToUploadFile
		}
		AlbumEntiy.Song = append(AlbumEntiy.Song, songservice.SongReqMapToSongEntity(SongRequestItem, SongResource, SongTypeArray, ArtistArray))
	}
	ArtistIdForAlbum := AlbumReq.Artist

	for _, ArtistItem := range ArtistIdForAlbum {
		AritstEntity, ErrorToGetArtist := ArtistRepo.GetArtistById(ArtistItem)
		if ErrorToGetArtist != nil {
			log.Print(ErrorToGetArtist)
			return MessageResponse{
				Message: "failed to get artist",
				Status:  "Failed",
			}, ErrorToGetArtist
		}
		AlbumEntiy.Artist = append(AlbumEntiy.Artist, AritstEntity)
	}
	ErrorToCompleteAlbum := AlbumRepo.UpdateAlbum(AlbumEntiy, AlbumId)
	if ErrorToCompleteAlbum != nil {
		log.Print(ErrorToCompleteAlbum)
		return MessageResponse{
			Message: "failed to Complete",
			Status:  "Failed",
		}, ErrorToCompleteAlbum
	}
	return MessageResponse{
		Message: "success to create album",
		Status:  "Success",
	}, nil

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
