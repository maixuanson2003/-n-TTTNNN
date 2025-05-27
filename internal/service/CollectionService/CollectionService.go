package collectionservice

import (
	"log"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
	"ten_module/internal/service/songservice"
	"time"
)

type CollectionService struct {
	SongRepo       *repository.SongRepository
	CollectionRepo *repository.CollectionRepostiory
}
type MessageResponse struct {
	Message string
	Status  string
}
type CollectionServiceInterface interface {
	GetListCollection() ([]response.CollectionResponse, error)
	GetCollectionById(Id int) (response.CollectionResponse, error)
	CreateCollection(NameCollection string) (MessageResponse, error)
	AddSongToCollection(SongId int, CollectionId int) (MessageResponse, error)
	DeleteCollection(SongId int, CollectionId int) (MessageResponse, error)
}

var CollectionServe *CollectionService

func InitCollectionService() {
	CollectionServe = &CollectionService{
		SongRepo:       repository.SongRepo,
		CollectionRepo: repository.CollectionRepo,
	}
}
func MapCollectToCollectionResponse(Collect entity.Collection) response.CollectionResponse {
	SongResponse := []response.SongResponse{}
	SongArray := Collect.Song
	for _, SongItem := range SongArray {
		SongResponse = append(SongResponse, songservice.SongEntityMapToSongResponse(SongItem))
	}
	return response.CollectionResponse{
		ID:             Collect.ID,
		NameCollection: Collect.NameCollection,
		CreateAt:       Collect.CreateAt,
		UpdateAt:       Collect.UpdateAt,
		Song:           SongResponse,
	}

}
func (CollectionServe *CollectionService) GetListCollection() ([]response.CollectionResponse, error) {
	CollectionRepo := CollectionServe.CollectionRepo
	CollectionArray, ErrorToGetCollect := CollectionRepo.FindAll()
	if ErrorToGetCollect != nil {
		log.Print(ErrorToGetCollect)
		return nil, ErrorToGetCollect
	}
	CollectionResponse := []response.CollectionResponse{}
	for _, CollectionItem := range CollectionArray {
		CollectionResponse = append(CollectionResponse, MapCollectToCollectionResponse(CollectionItem))
	}
	return CollectionResponse, nil

}
func (CollectionServe *CollectionService) GetCollectionById(Id int) (response.CollectionResponse, error) {
	CollectionRepo := CollectionServe.CollectionRepo
	Collection, ErrorToGetCollect := CollectionRepo.GetCollectById(Id)
	if ErrorToGetCollect != nil {
		log.Print(ErrorToGetCollect)
		return response.CollectionResponse{}, ErrorToGetCollect
	}
	CollectionResponse := MapCollectToCollectionResponse(Collection)
	return CollectionResponse, nil
}
func (CollectionServe *CollectionService) CreateCollection(NameCollection string) (MessageResponse, error) {
	CollectionRepo := CollectionServe.CollectionRepo
	Collect := entity.Collection{
		NameCollection: NameCollection,
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
	}
	ErrorToCreateCollection := CollectionRepo.CreateCollect(Collect)
	if ErrorToCreateCollection != nil {
		log.Print(ErrorToCreateCollection)
		return MessageResponse{
			Message: "failed to create Collection",
			Status:  "Failed",
		}, ErrorToCreateCollection
	}
	return MessageResponse{
		Message: "success to create Collection",
		Status:  "Success",
	}, nil

}
func (CollectionServe *CollectionService) AddSongToCollection(SongId int, CollectionId int) (MessageResponse, error) {
	SongRepo := CollectionServe.SongRepo
	CollectionRepo := CollectionServe.CollectionRepo
	SongItem, ErrorToGetSong := SongRepo.GetSongById(SongId)
	if ErrorToGetSong != nil {
		return MessageResponse{
			Message: "failed to add song",
			Status:  "Failed",
		}, ErrorToGetSong
	}
	CollectionItem, ErrorToGetCollect := CollectionRepo.GetCollectById(CollectionId)
	if ErrorToGetCollect != nil {
		return MessageResponse{
			Message: "failed to add song",
			Status:  "Failed",
		}, ErrorToGetCollect
	}
	CollectionItem.Song = append(CollectionItem.Song, SongItem)
	ErrorToAddSong := CollectionRepo.UpdateCollect(CollectionItem, CollectionId)
	if ErrorToAddSong != nil {
		return MessageResponse{
			Message: "failed to add song",
			Status:  "Failed",
		}, ErrorToAddSong
	}
	return MessageResponse{
		Message: "success to add song",
		Status:  "Success",
	}, nil
}
func (CollectionServe *CollectionService) DeleteSongFromCollection(SongId int, CollectionId int) (MessageResponse, error) {
	CollectionRepo := CollectionServe.CollectionRepo
	errorToDeleteSongFromRepo := CollectionRepo.DeleteSong(SongId, CollectionId)
	if errorToDeleteSongFromRepo != nil {
		log.Print(errorToDeleteSongFromRepo)
		return MessageResponse{
			Message: "faile to delete song",
			Status:  "Failed",
		}, errorToDeleteSongFromRepo
	}
	return MessageResponse{
		Message: "Success to delete song",
		Status:  "Success",
	}, nil
}
func (CollectionServe *CollectionService) DeleteCollectionByIds(CollectionId int) (MessageResponse, error) {
	CollectionRepo := CollectionServe.CollectionRepo
	errorToDelete := CollectionRepo.DeleteCollectById(CollectionId)
	if errorToDelete != nil {
		log.Print(errorToDelete)
		return MessageResponse{
			Message: "failed",
			Status:  "Failed",
		}, errorToDelete
	}
	return MessageResponse{
		Message: "success",
		Status:  "Success",
	}, nil
}
