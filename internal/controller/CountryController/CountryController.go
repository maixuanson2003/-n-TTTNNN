package countrycontroller

import (
	"encoding/json"
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
	r.HandleFunc("/update/country", CountryControll.UpdateCountry).Methods("PUT")

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
func (CountryControll *CountryController) UpdateCountry(Write http.ResponseWriter, Request *http.Request) {
	CountryName := Request.URL.Query().Get("namecountry")
	CountryId := Request.URL.Query().Get("countryid")
	CountryIdConvert, ErrorToConvert := strconv.Atoi(CountryId)
	if ErrorToConvert != nil {
		http.Error(Write, "failed to convert", http.StatusBadRequest)
		return
	}
	Resp, ErrorToUpdateCountry := CountryControll.CountryService.UpdateCountry(CountryName, CountryIdConvert)
	if ErrorToUpdateCountry != nil {
		http.Error(Write, "failed to update country", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)

}
