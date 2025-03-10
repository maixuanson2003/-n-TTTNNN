package artistservice

import (
	"log"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
)

type ArtistService struct {
	ArtistRepo *repository.ArtistRepository
}
type MessageResponse struct {
	Message string
	Status  string
}
type ArtistServiceInterface interface {
	GetListArtist() ([]response.ArtistResponse, error)
	CreateArtist(ArtistRequest request.ArtistRequest) (MessageResponse, error)
}

var ArtistServe *ArtistService

func InitArtistService() {
	ArtistServe = &ArtistService{
		ArtistRepo: repository.ArtistRepo,
	}
}
func MapArtistEntityToResponse(Artist entity.Artist) response.ArtistResponse {
	return response.ArtistResponse{
		ID:          Artist.ID,
		Name:        Artist.Name,
		BirthDay:    Artist.BirthDay,
		Description: Artist.Description,
		Country:     Artist.Country,
	}
}
func MapArtistToEntity(Artist request.ArtistRequest) entity.Artist {
	return entity.Artist{
		Name:        Artist.Name,
		BirthDay:    Artist.BirthDay,
		Description: Artist.Description,
		Country:     Artist.Country,
	}
}
func (ArtistServe *ArtistService) GetListArtist() ([]response.ArtistResponse, error) {
	ArtRepo := ArtistServe.ArtistRepo
	ArtistList, ErrorToGetArtist := ArtRepo.FindAll()
	ArtistRes := []response.ArtistResponse{}
	if ErrorToGetArtist != nil {
		log.Print(ErrorToGetArtist)
		return nil, ErrorToGetArtist
	}
	for _, ArtistItem := range ArtistList {
		ArtistRes = append(ArtistRes, MapArtistEntityToResponse(ArtistItem))
	}
	return ArtistRes, nil
}
func (ArtistServe *ArtistService) CreateArtist(ArtistRequest request.ArtistRequest) (MessageResponse, error) {
	ArtRepo := ArtistServe.ArtistRepo
	Artist := MapArtistToEntity(ArtistRequest)
	ErrorToCreateAritst := ArtRepo.CreateArtist(Artist)
	if ErrorToCreateAritst != nil {
		log.Print(ErrorToCreateAritst)
		return MessageResponse{
			Message: "Failed",
			Status:  "Failed",
		}, ErrorToCreateAritst
	}
	return MessageResponse{
		Message: "Success",
		Status:  "Success",
	}, nil
}
