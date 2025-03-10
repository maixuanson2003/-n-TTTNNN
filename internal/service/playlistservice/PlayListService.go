package playlistservice

import (
	"log"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
	"ten_module/internal/service/songservice"
	"time"
)

type PlayListService struct {
	PlayListRepo *repository.PlayListRepository
	UserRepo     *repository.UserRepository
	SongRepo     *repository.SongRepository
}
type MessageResponse struct {
	Message string
	Status  string
}
type PlayListServiceInterface interface {
	GetSongByPLayList(PlayListId int) ([]response.SongResponse, error)
	GetPlayListByUser(UserId string) ([]response.PlayListResponse, error)
	CreatePlayList(NamePlayList string, UserId string) (MessageResponse, error)
	AddSongToPlayList(SongId int, PlayListId int) (MessageResponse, error)
}

var PlayListServe *PlayListService

func InitPlayListService() {
	PlayListServe = &PlayListService{
		PlayListRepo: repository.PlayListRepo,
		UserRepo:     repository.UserRepo,
		SongRepo:     repository.SongRepo,
	}
}
func (playlistService *PlayListService) GetSongByPLayList(PlayListId int) ([]response.SongResponse, error) {
	PlayRepo := playlistService.PlayListRepo
	PlayList, ErrorToGetPlayList := PlayRepo.GetPlayListById(PlayListId)
	if ErrorToGetPlayList != nil {
		log.Print(ErrorToGetPlayList)
		return nil, ErrorToGetPlayList
	}
	SongArray := PlayList.Song
	ListSongResponse := []response.SongResponse{}
	for _, SongItem := range SongArray {
		ListSongResponse = append(ListSongResponse, songservice.SongEntityMapToSongResponse(SongItem))
	}

	return ListSongResponse, nil
}
func (playlistService *PlayListService) GetPlayListByUser(UserId string) ([]response.PlayListResponse, error) {
	UserRepo := playlistService.UserRepo
	User, ErrorToGetUser := UserRepo.FindById(UserId)
	if ErrorToGetUser != nil {
		log.Print(ErrorToGetUser)
		return nil, ErrorToGetUser
	}
	PlayList := User.PlayList
	PlayListResponse := []response.PlayListResponse{}
	for _, PlayListItem := range PlayList {
		PlayListResponse = append(PlayListResponse, response.PlayListResponse{
			ID:        PlayListItem.ID,
			Name:      PlayListItem.Name,
			CreateDay: PlayListItem.CreateDay,
		})

	}
	return PlayListResponse, nil
}
func (playlistService *PlayListService) CreatePlayList(NamePlayList string, UserId string) (MessageResponse, error) {
	PlayRepo := playlistService.PlayListRepo
	NewPlayList := entity.PlayList{
		Name:      NamePlayList,
		CreateDay: time.Now(),
		UpdateDay: time.Now(),
		UserId:    UserId,
	}
	ErrorToCreatePlayList := PlayRepo.CreatePlayList(NewPlayList)
	if ErrorToCreatePlayList != nil {
		log.Print(ErrorToCreatePlayList)
		return MessageResponse{
			Message: "failed to create playlist",
			Status:  "Failed",
		}, ErrorToCreatePlayList

	}
	return MessageResponse{
		Message: "success to create playlist",
		Status:  "Success",
	}, nil
}
func (playlistService *PlayListService) AddSongToPlayList(SongId int, PlayListId int) (MessageResponse, error) {
	SongRepo := playlistService.SongRepo
	PlayListRepo := playlistService.PlayListRepo
	SongItem, ErrorToGetSong := SongRepo.GetSongById(SongId)
	PlayList, ErrorToGetPlayList := PlayListRepo.GetPlayListById(PlayListId)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return MessageResponse{
			Message: "failed to add song playlist",
			Status:  "Failed",
		}, ErrorToGetSong

	}
	if ErrorToGetPlayList != nil {
		log.Print(ErrorToGetPlayList)
		return MessageResponse{
			Message: "failed to add song playlist",
			Status:  "Failed",
		}, ErrorToGetPlayList
	}
	PlayList.Song = append(PlayList.Song, SongItem)
	ErrorToUpdatePlayList := PlayListRepo.UpdatePlayList(PlayList, PlayListId)
	if ErrorToUpdatePlayList != nil {
		log.Print(ErrorToUpdatePlayList)
		return MessageResponse{
			Message: "failed to add song playlist",
			Status:  "Failed",
		}, ErrorToUpdatePlayList
	}
	return MessageResponse{
		Message: "success to add song playlist",
		Status:  "Success",
	}, nil
}
