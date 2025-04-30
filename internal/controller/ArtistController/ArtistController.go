package artistcontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"ten_module/internal/DTO/request"
	"ten_module/internal/service/artistservice"

	"github.com/gorilla/mux"
)

type ArtistController struct {
	ArtistService *artistservice.ArtistService
}

var ArtistControll *ArtistController

func InitArtistController() {
	ArtistControll = &ArtistController{
		ArtistService: artistservice.ArtistServe,
	}
}
func (ArtController *ArtistController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/listart", ArtController.GetListArtist).Methods("GET")
	r.HandleFunc("/createart", ArtController.CreateArtist).Methods("POST")
	r.HandleFunc("/searchart", ArtController.SearchArtist).Methods("GET")
	r.HandleFunc("/filterart", ArtController.FilterArtist).Methods("GET")
	r.HandleFunc("/createindex", ArtController.CreateAritstIndexToElastic).Methods("POST")
	r.HandleFunc("/addart", ArtController.AddAritstToElastic).Methods("POST")
	r.HandleFunc("/artist/{id}", ArtController.GetArtistById).Methods("GET")
}
func (ArtController *ArtistController) GetListArtist(write http.ResponseWriter, Request *http.Request) {
	Artist, ErrorToGetList := ArtController.ArtistService.GetListArtist()
	if ErrorToGetList != nil {
		log.Print(ErrorToGetList)
		http.Error(write, "Fail to get List", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(write).Encode(Artist)
}
func (ArtController *ArtistController) GetArtistById(write http.ResponseWriter, Request *http.Request) {
	url := Request.URL.Path
	AritstIdParam := strings.Split(url, "/")[3]
	ArtistId, errorToConvert := strconv.Atoi(AritstIdParam)
	if errorToConvert != nil {
		log.Print(errorToConvert)
		http.Error(write, "Fail to get artist", http.StatusBadRequest)
		return
	}
	Artist, ErrorToGetList := ArtController.ArtistService.GetArtistById(ArtistId)
	if ErrorToGetList != nil {
		log.Print(ErrorToGetList)
		http.Error(write, "Fail to get List", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(write).Encode(Artist)
}
func (ArtController *ArtistController) CreateArtist(write http.ResponseWriter, Request *http.Request) {
	var artistRequest request.ArtistRequest
	json.NewDecoder(Request.Body).Decode(&artistRequest)
	Resp, ErrorToCreateArtist := ArtController.ArtistService.CreateArtist(artistRequest)
	if ErrorToCreateArtist != nil {
		log.Print(ErrorToCreateArtist)
		http.Error(write, "Failed To CreateSong", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(write).Encode(Resp)
}
func (ArtController *ArtistController) SearchArtist(write http.ResponseWriter, Request *http.Request) {
	KeyWord := Request.URL.Query().Get("keyword")
	Resp, errorToSearchArtist := ArtController.ArtistService.SearchArtist(KeyWord)
	if errorToSearchArtist != nil {
		http.Error(write, "failed to search", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(write).Encode(Resp)

}
func (ArtController *ArtistController) FilterArtist(write http.ResponseWriter, Request *http.Request) {
	CountryIdQuery := Request.URL.Query().Get("countryid")
	CountryId, errorToConvert := strconv.Atoi(CountryIdQuery)
	if errorToConvert != nil {
		http.Error(write, "failed to search", http.StatusBadRequest)
		return
	}

	Resp, errorToSearchArtist := ArtController.ArtistService.FilterArtist(CountryId)
	if errorToSearchArtist != nil {
		http.Error(write, "failed to search", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(write).Encode(Resp)

}
func (ArtController *ArtistController) CreateAritstIndexToElastic(write http.ResponseWriter, Request *http.Request) {
	ArtController.ArtistService.CreateIndexArtistInElastic()
	write.Header().Set("Content-Type", "text/html")
	write.WriteHeader(http.StatusAccepted)
	write.Write([]byte(`{"message": "Index created successfully"}`))
}
func (ArtController *ArtistController) AddAritstToElastic(write http.ResponseWriter, Request *http.Request) {
	errors := ArtController.ArtistService.AddArtistToElastic()
	if errors != nil {
		http.Error(write, "faile to add", http.StatusBadRequest)
		return
	}
	write.Header().Set("Content-Type", "text/html")
	write.WriteHeader(http.StatusAccepted)
	write.Write([]byte(`{"message": "add succes"}`))
}
