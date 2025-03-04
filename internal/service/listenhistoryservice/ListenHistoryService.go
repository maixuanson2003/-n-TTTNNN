package listenhistoryservice

import (
	"log"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
	"time"
)

type ListenHistoryService struct {
	HistoryRepo *repository.ListenHistoryRepo
}

var HistoryService *ListenHistoryService

func InitListenHistoryService() {
	HistoryService = &ListenHistoryService{
		HistoryRepo: repository.ListenRepo,
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
	History := entity.ListenHistory{
		SongId:    SongId,
		UserId:    UserId,
		ListenDay: time.Now(),
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
