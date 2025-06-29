package artistservice

import (
	"log"
	"mime/multipart"
	"ten_module/internal/Config"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/Helper/elastichelper"
	"ten_module/internal/repository"
	"ten_module/internal/service/albumservice"
	"ten_module/internal/service/songservice"
)

type ArtistService struct {
	ArtistRepo *repository.ArtistRepository
	CountryRep *repository.CountryRepository
	SongRep    *repository.SongRepository
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
	GetArtistById(artistId int) ([]map[string]interface{}, error)
}

var ArtistServe *ArtistService

func InitArtistService() {
	ArtistServe = &ArtistService{
		ArtistRepo: repository.ArtistRepo,
		CountryRep: repository.CountryRepo,
		SongRep:    repository.SongRepo,
	}
}
func MapArtistEntityToResponse(Artist entity.Artist, NameCountry string) response.ArtistResponse {
	return response.ArtistResponse{
		ID:          Artist.ID,
		Name:        Artist.Name,
		BirthDay:    Artist.BirthDay,
		Description: Artist.Description,
		Country:     NameCountry,
		Image:       Artist.Image,
	}
}
func MapArtistToEntity(Artist request.ArtistRequest, Image string) entity.Artist {
	return entity.Artist{
		Name:        Artist.Name,
		BirthDay:    Artist.BirthDay,
		Description: Artist.Description,
		CountryId:   Artist.CountryId,
		Image:       Image,
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
func (ArtistServe *ArtistService) CreateArtist(ArtistRequest request.ArtistRequest, File *multipart.FileHeader) (MessageResponse, error) {
	ArtRepo := ArtistServe.ArtistRepo
	go func() {
		errs := ArtistServe.processArtistBackground(*ArtRepo, ArtistRequest, File)
		if errs != nil {
			log.Print(errs)
		}
	}()
	return MessageResponse{
		Message: "Success",
		Status:  "Success",
	}, nil
}
func (ArtistServe *ArtistService) processArtistBackground(ArtRepo repository.ArtistRepository, ArtistRequest request.ArtistRequest, File *multipart.FileHeader) error {
	imageResource, errs := Config.HandleUploadImage(File, ArtistRequest.Name)
	if errs != nil {
		return errs
	}
	Artist := MapArtistToEntity(ArtistRequest, imageResource)
	ErrorToCreateAritst := ArtRepo.CreateArtist(Artist)
	return ErrorToCreateAritst
}
func (ArtistServe *ArtistService) processArtistUpdate(ArtRepo repository.ArtistRepository, artist entity.Artist, File *multipart.FileHeader) error {
	imageResource, errs := Config.HandleUploadImage(File, artist.Name)
	if errs != nil {
		return errs
	}
	artist.Image = imageResource
	ErrorToCreateAritst := ArtRepo.UpdateAritst(artist, artist.ID)
	return ErrorToCreateAritst
}
func (ArtistServe *ArtistService) UpdateArtist(ArtistRequest request.ArtistRequest, artistId int, File *multipart.FileHeader) (MessageResponse, error) {
	ArtRepo := ArtistServe.ArtistRepo
	Artist, err := ArtRepo.GetArtistById(artistId)
	if err != nil {
		log.Print(err)
		return MessageResponse{}, err
	}
	Artist.Name = ArtistRequest.Name
	Artist.BirthDay = ArtistRequest.BirthDay
	Artist.CountryId = ArtistRequest.CountryId
	Artist.Description = ArtistRequest.Description
	ErrorToCreateAritst := ArtRepo.UpdateAritst(Artist, artistId)
	if ErrorToCreateAritst != nil {
		log.Print(ErrorToCreateAritst)
		return MessageResponse{
			Message: "Failed",
			Status:  "Failed",
		}, ErrorToCreateAritst
	}
	go func() {
		errs := ArtistServe.processArtistUpdate(*ArtRepo, Artist, File)
		if errs != nil {
			log.Print(errs)
		}
	}()
	return MessageResponse{
		Message: "Success",
		Status:  "Success",
	}, nil
}

// func (ArtistServe *ArtistService) SearchArtistByKeyWord(Keyword string) ([]response.ArtistResponse, error) {
// 	ArtistRepo := ArtistServe.ArtistRepo
// 	CountryRepo := ArtistServe.CountryRep
// 	ArtistList, errorToSearch := ArtistRepo.SearchArtist(Keyword)
// 	if errorToSearch != nil {
// 		log.Print(errorToSearch)
// 		return nil, errorToSearch
// 	}
// 	ArtitsResponse := []response.ArtistResponse{}
// 	for _, Item := range ArtistList {
// 		Country, errorToGetCountry := CountryRepo.GetCountryById(Item.CountryId)
// 		if errorToGetCountry != nil {
// 			log.Print(errorToGetCountry)
// 			return nil, errorToGetCountry
// 		}
// 		ArtitsResponse = append(ArtitsResponse, MapArtistEntityToResponse(Item, Country.CountryName))

// 	}

//		return ArtitsResponse, nil
//	}
func (ArtistServe *ArtistService) SearchArtist(Keyword string) ([]response.ArtistResponse, error) {
	ArtistRepo := ArtistServe.ArtistRepo
	CountryRepo := ArtistServe.CountryRep
	ArtistList, errorToSearch := ArtistRepo.SearchArtist(Keyword)
	if errorToSearch != nil {
		log.Print(errorToSearch)
		return nil, errorToSearch
	}
	ArtitsResponse := []response.ArtistResponse{}
	for _, Item := range ArtistList {
		Country, errorToGetCountry := CountryRepo.GetCountryById(Item.CountryId)
		if errorToGetCountry != nil {
			log.Print(errorToGetCountry)
			return nil, errorToGetCountry
		}
		ArtitsResponse = append(ArtitsResponse, MapArtistEntityToResponse(Item, Country.CountryName))

	}

	return ArtitsResponse, nil

}
func (ArtistServe *ArtistService) FilterArtist(CountryId int) ([]response.ArtistResponse, error) {
	ArtistRepo := ArtistServe.ArtistRepo
	CountryRepo := ArtistServe.CountryRep
	ArtistList, errorToSearch := ArtistRepo.FilterArtist(CountryId)
	if errorToSearch != nil {
		log.Print(errorToSearch)
		return nil, errorToSearch
	}
	ArtitsResponse := []response.ArtistResponse{}
	for _, Item := range ArtistList {
		Country, errorToGetCountry := CountryRepo.GetCountryById(Item.CountryId)
		if errorToGetCountry != nil {
			log.Print(errorToGetCountry)
			return nil, errorToGetCountry
		}
		ArtitsResponse = append(ArtitsResponse, MapArtistEntityToResponse(Item, Country.CountryName))

	}

	return ArtitsResponse, nil

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
func (ArtistServe *ArtistService) GetArtistById(artistId int) (map[string]interface{}, error) {
	ArtistRepo := ArtistServe.ArtistRepo
	CountryRepo := ArtistServe.CountryRep
	Artist, errorToGetArtist := ArtistRepo.GetArtistById(artistId)
	if errorToGetArtist != nil {
		return nil, errorToGetArtist
	}
	Country, errorToGetCountry := CountryRepo.GetCountryById(Artist.CountryId)
	if errorToGetCountry != nil {
		return nil, errorToGetCountry
	}
	SongList := Artist.Song
	AlbumList := Artist.Album
	ArtistResponse := MapArtistEntityToResponse(Artist, Country.CountryName)
	SongResponse := []response.SongResponse{}
	AlbumResponse := []response.AlbumResponse{}
	for _, SongItem := range SongList {
		SongResponse = append(SongResponse, songservice.SongEntityMapToSongResponse(SongItem))

	}
	for _, AlbumItem := range AlbumList {
		AlbumResponse = append(AlbumResponse, albumservice.AlbumEntityMapToAlbumResponse(AlbumItem, ArtistServe.CountryRep, ArtistServe.SongRep))
	}
	response := map[string]interface{}{
		"artist": ArtistResponse,
		"song":   SongResponse,
		"album":  AlbumResponse,
	}
	return response, nil
}
func (ArtistServe *ArtistService) DeleteArtist(artistId int) (MessageResponse, error) {
	ArtistRepo := ArtistServe.ArtistRepo
	err := ArtistRepo.DeleteArtist(artistId)
	if err != nil {
		return MessageResponse{
			Message: "failed",
			Status:  "Failed",
		}, err
	}
	return MessageResponse{
		Message: "Success",
		Status:  "Success",
	}, nil
}
