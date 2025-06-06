package countryservice

import (
	"log"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
	"time"
)

type CountryService struct {
	CountryRepo *repository.CountryRepository
}
type MessageResponse struct {
	Message string
	Status  string
}

var CountryServe *CountryService

func InitCountryService() {
	CountryServe = &CountryService{
		CountryRepo: repository.CountryRepo,
	}
}

type CountryServiceInterface interface {
	GetListCountry() ([]response.CountryResponse, error)
	CreateCountry(CountryName string) (MessageResponse, error)
	UpdateCountry(CountryName string, Id int) (MessageResponse, error)
}

func (CountryServe *CountryService) GetListCountry() ([]response.CountryResponse, error) {
	CountryRepo := CountryServe.CountryRepo
	CountryList, ErrorToGetList := CountryRepo.FindAll()
	if ErrorToGetList != nil {
		log.Print(ErrorToGetList)
		return nil, ErrorToGetList
	}
	CountryResponse := []response.CountryResponse{}
	for _, CountryItem := range CountryList {
		CountryResponse = append(CountryResponse, response.CountryResponse{
			Id:          CountryItem.ID,
			CountryName: CountryItem.CountryName,
			CreateAt:    CountryItem.CreateAt,
		})
	}
	return CountryResponse, nil

}
func (CountryServe *CountryService) CreateCountry(CountryName string) (MessageResponse, error) {
	CountryRepo := CountryServe.CountryRepo
	NewCountry := entity.Country{
		CountryName: CountryName,
		CreateAt:    time.Now(),
	}
	ErrorToCreateCountry := CountryRepo.CreateCountry(NewCountry)
	if ErrorToCreateCountry != nil {
		log.Print(ErrorToCreateCountry)
		return MessageResponse{
			Message: "failed to create country",
			Status:  "Failed",
		}, ErrorToCreateCountry
	}
	return MessageResponse{
		Message: "success to create country",
		Status:  "Success",
	}, nil

}
func (CountryServe *CountryService) UpdateCountry(CountryName string, Id int) (MessageResponse, error) {
	CountryRepo := CountryServe.CountryRepo
	Country, ErrorToGetCountry := CountryRepo.GetCountryById(Id)
	if ErrorToGetCountry != nil {
		log.Print(ErrorToGetCountry)
		return MessageResponse{
			Message: "faile to update country",
			Status:  "Failed",
		}, ErrorToGetCountry
	}
	Country.CountryName = CountryName
	ErrorToUpdate := CountryRepo.UpdateCountry(Country, Id)
	if ErrorToUpdate != nil {
		log.Print(ErrorToUpdate)
		return MessageResponse{
			Message: "faile to update country",
			Status:  "Failed",
		}, ErrorToUpdate
	}
	return MessageResponse{
		Message: "success to update country",
		Status:  "Success",
	}, nil

}
func (CountryServe *CountryService) DeleteCountryById(Id int) (MessageResponse, error) {
	CountryRepo := CountryServe.CountryRepo
	_, err := CountryRepo.GetCountryById(Id)
	if err != nil {
		log.Print(err)
		return MessageResponse{
			Message: "failed to delete country, not found",
			Status:  "Failed",
		}, err
	}

	// Thực hiện xoá
	err = CountryRepo.DeleteCountryByID(Id)
	if err != nil {
		log.Print(err)
		return MessageResponse{
			Message: "failed to delete country",
			Status:  "Failed",
		}, err
	}

	return MessageResponse{
		Message: "success to delete country",
		Status:  "Success",
	}, nil
}
func (CountryServe *CountryService) GetCountryById(Id int) (response.CountryResponse, error) {
	CountryRepo := CountryServe.CountryRepo

	CountryEntity, err := CountryRepo.GetCountryById(Id)
	if err != nil {
		log.Print(err)
		return response.CountryResponse{}, err
	}

	CountryResp := response.CountryResponse{
		Id:          CountryEntity.ID,
		CountryName: CountryEntity.CountryName,
		CreateAt:    CountryEntity.CreateAt,
	}

	return CountryResp, nil
}
