package songtypeservice

import (
	"log"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
	"time"
)

type SongTypeService struct {
	SongTypeRepo *repository.SongTypeRepository
}
type MessageResponse struct {
	Message string
	Status  string
}
type SongTypeServiceInterface interface {
	GetListSongType() ([]map[string]interface{}, error)
	CreateSongType(Type string) (MessageResponse, error)
	UpdateSongType(Type string, Id int) (MessageResponse, error)
	DeleteSongTypeById(Id int) (MessageResponse, error)
}

var SongTypeServe *SongTypeService

func InitSongTypeService() {
	SongTypeServe = &SongTypeService{
		SongTypeRepo: repository.SongTypeRepo,
	}
}

func (SongTypeServe *SongTypeService) GetListSongType() ([]map[string]interface{}, error) {
	SongTypeRepo := SongTypeServe.SongTypeRepo
	SongTypeList, errorToGetList := SongTypeRepo.FindAll()
	if errorToGetList != nil {
		log.Print(errorToGetList)
		return nil, errorToGetList
	}
	SongTypeResponse := []map[string]interface{}{}
	for _, SongTypeItem := range SongTypeList {
		ResponseItem := map[string]interface{}{
			"id":     SongTypeItem.ID,
			"type":   SongTypeItem.Type,
			"create": SongTypeItem.CreatAt,
		}
		SongTypeResponse = append(SongTypeResponse, ResponseItem)

	}
	return SongTypeResponse, nil

}
func (SongTypeServe *SongTypeService) CreateSongType(Type string) (MessageResponse, error) {
	SongTypeRepo := SongTypeServe.SongTypeRepo
	SongType := entity.SongType{
		Type:    Type,
		CreatAt: time.Now(),
	}
	errorToCreateSongType := SongTypeRepo.CreateSongType(SongType)
	if errorToCreateSongType != nil {
		return MessageResponse{
			Message: "failed to create",
			Status:  "Failed",
		}, errorToCreateSongType
	}
	return MessageResponse{
		Message: "success to create",
		Status:  "Success",
	}, nil
}
func (SongTypeServe *SongTypeService) UpdateSongType(Type string, Id int) (MessageResponse, error) {
	SongTypeRepo := SongTypeServe.SongTypeRepo
	SongTypeItem, errorToGetSongType := SongTypeRepo.GetSongTypeById(Id)
	if errorToGetSongType != nil {
		return MessageResponse{
			Message: "failed to update",
			Status:  "Failed",
		}, errorToGetSongType
	}
	SongTypeItem.Type = Type
	errorToUpdateSongType := SongTypeRepo.UpdateSongType(SongTypeItem, Id)
	if errorToUpdateSongType != nil {
		return MessageResponse{
			Message: "failed to update",
			Status:  "Failed",
		}, errorToGetSongType
	}
	return MessageResponse{
		Message: "success to update",
		Status:  "Success",
	}, nil
}
