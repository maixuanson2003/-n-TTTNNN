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
	}
}

type AlbumServiceInterface interface {
	GetListAlbum() ([]response.AlbumResponse, error)
	GetAlbumById(Id int) (response.AlbumResponse, error)
	CreateAlbum(AlbumReq request.AlbumRequest, SongFileAlum []request.SongFileAlbum) (MessageResponse, error)
	UpdateAlbum(AlbumReq request.AlbumRequest) (MessageResponse, error)
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
