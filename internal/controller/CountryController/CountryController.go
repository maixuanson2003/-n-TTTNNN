package countrycontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"ten_module/internal/service/countryservice"

	"github.com/gorilla/mux"
)

type CountryController struct {
	CountryService *countryservice.CountryService
}

var CountryControll *CountryController

func InitCountryController() {
	CountryControll = &CountryController{
		CountryService: countryservice.CountryServe,
	}
}
func (CountryControll *CountryController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/country/list", CountryControll.GetListCountry).Methods("GET")
	r.HandleFunc("/create/country", CountryControll.CreateCountry).Methods("POST")
	r.HandleFunc("/updateregion", CountryControll.UpdateCountryRegion).Methods("PUT")
	r.HandleFunc("/country/{id}", CountryControll.GetCountryById).Methods("GET")
	r.HandleFunc("/deletecountry/{id}", CountryControll.DeleteCountryById).Methods("DELETE")
}
func (CountryControll *CountryController) GetListCountry(Write http.ResponseWriter, Request *http.Request) {
	Resp, ErrorToGetList := CountryControll.CountryService.GetListCountry()
	print(Resp)
	if ErrorToGetList != nil {
		http.Error(Write, "failed to get list", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
func (CountryControll *CountryController) CreateCountry(Write http.ResponseWriter, Request *http.Request) {
	CountryName := Request.URL.Query().Get("namecountry")
	Resp, ErrorToCreateCountry := CountryControll.CountryService.CreateCountry(CountryName)
	if ErrorToCreateCountry != nil {
		http.Error(Write, "failed to create Country", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)

}
func (CountryControll *CountryController) UpdateCountryRegion(Write http.ResponseWriter, Request *http.Request) {
	CountryName := Request.URL.Query().Get("namecountry")
	CountryId := Request.URL.Query().Get("countryid")
	CountryIdConvert, ErrorToConvert := strconv.Atoi(CountryId)
	if ErrorToConvert != nil {
		log.Print(ErrorToConvert)
		http.Error(Write, "failed to convert", http.StatusBadRequest)
		return
	}
	Resp, ErrorToUpdateCountry := CountryControll.CountryService.UpdateCountry(CountryName, CountryIdConvert)
	if ErrorToUpdateCountry != nil {
		log.Print(ErrorToUpdateCountry)
		http.Error(Write, "failed to update country", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)

}
func (CountryControll *CountryController) GetCountryById(Write http.ResponseWriter, Request *http.Request) {
	params := mux.Vars(Request)
	idParam := params["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(Write, "invalid country id", http.StatusBadRequest)
		return
	}

	resp, err := CountryControll.CountryService.GetCountryById(id)
	if err != nil {
		http.Error(Write, "country not found", http.StatusNotFound)
		return
	}

	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(resp)
}
func (CountryControll *CountryController) DeleteCountryById(Write http.ResponseWriter, Request *http.Request) {
	params := mux.Vars(Request)
	idParam := params["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Print(err)
		http.Error(Write, "invalid country id", http.StatusBadRequest)
		return
	}

	resp, err := CountryControll.CountryService.DeleteCountryById(id)
	if err != nil {
		log.Print(err)
		http.Error(Write, "failed to delete country", http.StatusInternalServerError)
		return
	}

	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(resp)
}
