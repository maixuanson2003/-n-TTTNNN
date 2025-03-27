package artistservice

import (
	"log"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/Helper/elastichelper"
	"ten_module/internal/repository"
)

type ArtistService struct {
	ArtistRepo *repository.ArtistRepository
	CountryRep *repository.CountryRepository
}
type MessageResponse struct {
	Message string
	Status  string
}
type ArtistServiceInterface interface {
	GetListArtist() ([]response.ArtistResponse, error)
	CreateArtist(ArtistRequest request.ArtistRequest) (MessageResponse, error)
	SearchArtist(Keyword string) ([]response.ArtistResponse, error)
	AddArtistToElastic() error
	CreateIndexArtistInElastic()
}

var ArtistServe *ArtistService

func InitArtistService() {
	ArtistServe = &ArtistService{
		ArtistRepo: repository.ArtistRepo,
		CountryRep: repository.CountryRepo,
	}
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
func MapArtistToEntity(Artist request.ArtistRequest) entity.Artist {
	return entity.Artist{
		Name:        Artist.Name,
		BirthDay:    Artist.BirthDay,
		Description: Artist.Description,
		CountryId:   Artist.CountryId,
	}
}
func (ArtistServe *ArtistService) GetListArtist() ([]response.ArtistResponse, error) {
	ArtRepo := ArtistServe.ArtistRepo
	CountryRepo := ArtistServe.CountryRep
	ArtistList, ErrorToGetArtist := ArtRepo.FindAll()
	ArtistRes := []response.ArtistResponse{}
	if ErrorToGetArtist != nil {
		log.Print(ErrorToGetArtist)
		return nil, ErrorToGetArtist
	}
	for _, ArtistItem := range ArtistList {
		Country, ErrorToGetCountry := CountryRepo.GetCountryById(ArtistItem.CountryId)
		if ErrorToGetCountry != nil {
			log.Print(ErrorToGetCountry)
			return nil, ErrorToGetCountry
		}
		ArtistRes = append(ArtistRes, MapArtistEntityToResponse(ArtistItem, Country.CountryName))
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
func (ArtistServe *ArtistService) SearchArtist(Keyword string) ([]response.ArtistResponse, error) {
	Elastic := elastichelper.ElasticHelpers
	ArtistResponse, errorToSearch := Elastic.SearchArtist(Keyword)
	if errorToSearch != nil {
		log.Print(errorToSearch)
		return nil, errorToSearch
	}
	return ArtistResponse, nil

}
func (ArtistServe *ArtistService) AddArtistToElastic() error {
	ArtistRepo := ArtistServe.ArtistRepo
	ElasticHelper := elastichelper.ElasticHelpers
	ArtistList, ErrorToGetArtist := ArtistRepo.FindAll()
	if ErrorToGetArtist != nil {
		log.Print(ErrorToGetArtist)
		return ErrorToGetArtist
	}
	errors := ElasticHelper.InsertDataArtistToIndex("artist", ArtistList)
	if errors != nil {
		log.Print(errors)
		return errors
	}
	return nil
}
func (ArtistServe *ArtistService) CreateIndexArtistInElastic() {
	ElasticHelper := elastichelper.ElasticHelpers
	errors := ElasticHelper.CreateIndexElastic("artist")
	if errors != nil {
		log.Print("da co index")
		return
	}
}
