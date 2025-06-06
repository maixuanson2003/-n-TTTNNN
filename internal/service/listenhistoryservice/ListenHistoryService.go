package listenhistoryservice

import (
	"log"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
	"time"
)

type ListenHistoryService struct {
	HistoryRepo *repository.ListenHistoryRepo
	SongRepo    *repository.SongRepository
}

var HistoryService *ListenHistoryService

func InitListenHistoryService() {
	HistoryService = &ListenHistoryService{
		HistoryRepo: repository.ListenRepo,
		SongRepo:    repository.SongRepo,
	}
}

type MessageResponse struct {
	Message string
	Status  string
}
type ListenHistoryServiceInterface interface {
	SaveHistoryListen(UserId string, SongId int) (MessageResponse, error)
}

func (HistoryServ *ListenHistoryService) SaveHistoryListen(UserId string, SongId int) (MessageResponse, error) {
	HistoryRepo := HistoryServ.HistoryRepo
	SongRepo := HistoryServ.SongRepo
	log.Print(UserId)
	var userIdPtr *string
	if UserId != "" {
		userIdPtr = &UserId
	}
	History := entity.ListenHistory{
		SongId:    SongId,
		UserId:    userIdPtr,
		ListenDay: time.Now(),
	}
	SongItem, ErrorToGetSong := SongRepo.GetSongById(SongId)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return MessageResponse{
			Message: "failed to get song",
			Status:  "Failed",
		}, ErrorToGetSong

	}
	SongItem.ListenAmout += 1
	ErrorToUpdateSong := SongRepo.UpdateSong(SongItem, SongId)
	if ErrorToUpdateSong != nil {
		return MessageResponse{
			Message: "failed to update song",
			Status:  "Failed",
		}, ErrorToUpdateSong
	}
	ErrorToSaveHistory := HistoryRepo.CreateHistory(History)
	if ErrorToSaveHistory != nil {
		log.Print(ErrorToSaveHistory)
		return MessageResponse{
			Message: "failed to save history",
			Status:  "Failed",
		}, ErrorToSaveHistory
	}
	return MessageResponse{
		Message: "success to save history",
		Status:  "Success",
	}, nil

}
