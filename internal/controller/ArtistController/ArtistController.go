package artistcontroller

import (
	"encoding/json"
	"log"
	"net/http"
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
