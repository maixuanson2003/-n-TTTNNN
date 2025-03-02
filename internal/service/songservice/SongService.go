package songservice

import (
	"errors"
	"log"
	"net/http"
	"ten_module/internal/Config"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
	"time"
)

type SongService struct {
	UserRepo     *repository.UserRepository
	SongRepo     *repository.SongRepository
	SongTypeRepo *repository.SongTypeRepository
	ArtistRepo   *repository.ArtistRepository
}
type SongServiceInterface interface {
	GetSongById(Id int) (response.SongResponse, error)
	GetAllSong() ([]response.SongResponse, error)
	CreateNewSong(SongReq request.SongRequest, SongFile request.SongFile) (MessageResponse, error)
	DownLoadSong(Id int) (SongDownload, error)
}
type MessageResponse struct {
	Message string
	Status  string
}

var SongServices *SongService

func InitSongService() {
	SongServices = &SongService{
		UserRepo:     repository.UserRepo,
		SongRepo:     repository.SongRepo,
		SongTypeRepo: repository.SongTypeRepo,
		ArtistRepo:   repository.ArtistRepo,
	}
}
func SongReqMapToSongEntity(SongReq request.SongRequest, resource string, ListSongType []entity.SongType, ListArtist []entity.Artist) entity.Song {
	return entity.Song{
		NameSong:     SongReq.NameSong,
		Description:  SongReq.Description,
		ReleaseDay:   time.Now(),
		CreateDay:    time.Now(),
		UpdateDay:    time.Now(),
		Point:        SongReq.Point,
		LikeAmount:   0,
		Status:       "Release",
		CountryId:    SongReq.CountryId,
		ListenAmout:  0,
		SongResource: resource,
		SongType:     ListSongType,
		Artist:       ListArtist,
	}
}
func SongEntityMapToSongResponse(Song entity.Song) response.SongResponse {
	return response.SongResponse{
		ID:           Song.ID,
		NameSong:     Song.NameSong,
		Description:  Song.Description,
		ReleaseDay:   Song.ReleaseDay,
		CreateDay:    Song.CreateDay,
		UpdateDay:    Song.UpdateDay,
		Point:        Song.Point,
		LikeAmount:   Song.LikeAmount,
		Status:       Song.Status,
		CountryId:    Song.CountryId,
		ListenAmout:  Song.ListenAmout,
		AlbumId:      Song.AlbumId,
		SongResource: Song.SongResource,
	}

}
func (songServe *SongService) CreateNewSong(SongReq request.SongRequest, SongFile request.SongFile) (MessageResponse, error) {
	ListSongType := []entity.SongType{}
	ListArtist := []entity.Artist{}
	for _, IdSongType := range SongReq.SongType {
		SongType, err := songServe.SongTypeRepo.GetSongTypeById(IdSongType)
		if err != nil {
			log.Print(err)
			return MessageResponse{}, err
		}
		ListSongType = append(ListSongType, SongType)
	}
	for _, IdArtist := range SongReq.Artist {
		Artist, err := songServe.ArtistRepo.GetArtistById(IdArtist)
		if err != nil {
			log.Print(err)
			return MessageResponse{}, err
		}
		ListArtist = append(ListArtist, Artist)
	}
	resourceSong, err := Config.HandleUpLoadFile(SongFile.File, SongReq.NameSong)
	if SongReq.NameSong == "" {
		return MessageResponse{
			Message: "Failed to create",
			Status:  "Failed",
		}, errors.New("name song is empty")
	}
	if err != nil {
		return MessageResponse{
			Message: "Failed to create",
			Status:  "Failed",
		}, err
	}
	SongEntity := SongReqMapToSongEntity(SongReq, resourceSong, ListSongType, ListArtist)
	errorToCreateSong := songServe.SongRepo.CreateSong(SongEntity)
	if errorToCreateSong != nil {
		return MessageResponse{
			Message: "failed to create song",
			Status:  "failed",
		}, errorToCreateSong
	}
	return MessageResponse{
		Message: "Success to create song",
		Status:  "Success",
	}, nil

}
func (songServe *SongService) GetAllSong() ([]response.SongResponse, error) {
	SongRepos := songServe.SongRepo
	ListSong, ErrorToGetListSong := SongRepos.FindAll()
	if ErrorToGetListSong != nil {
		log.Print(ErrorToGetListSong)
		return nil, ErrorToGetListSong
	}
	ListSongResponse := []response.SongResponse{}
	for _, SongItem := range ListSong {
		ListSongResponse = append(ListSongResponse, SongEntityMapToSongResponse(SongItem))
	}
	return ListSongResponse, nil
}
func (songServe *SongService) GetSongById(Id int) (response.SongResponse, error) {
	SongRepos := songServe.SongRepo
	Song, ErrorToGetSong := SongRepos.GetSongById(Id)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return response.SongResponse{}, ErrorToGetSong
	}
	SongResponse := SongEntityMapToSongResponse(Song)
	return SongResponse, nil
}

type SongDownload struct {
	Resp     *http.Response
	NameSong string
}

func (songServe *SongService) DownLoadSong(Id int) (SongDownload, error) {
	SongRepos := songServe.SongRepo
	Song, ErrorToGetSong := SongRepos.GetSongById(Id)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return SongDownload{}, ErrorToGetSong
	}
	resp, errorToGetSongAudio := Config.HandleDownLoadFile(Song.NameSong, "video")
	if errorToGetSongAudio != nil {
		log.Print(errorToGetSongAudio)
		return SongDownload{}, errorToGetSongAudio
	}
	return SongDownload{
		Resp:     resp,
		NameSong: Song.NameSong,
	}, nil
}
