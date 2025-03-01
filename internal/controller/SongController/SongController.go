package songcontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	middleware "ten_module/Middleware"
	"ten_module/internal/DTO/request"
	"ten_module/internal/service/songservice"

	"github.com/gorilla/mux"
)

type SongController struct {
	songService *songservice.SongService
	middlware   *middleware.UseMiddleware
}

var SongControllers *SongController

func InitSongController() {
	SongControllers = &SongController{
		songService: songservice.SongServices,
		middlware:   middleware.Middlewares,
	}
}
func (Controller *SongController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/song/create", Controller.CreateNewSong).Methods("POST")
}
func (Controller *SongController) CreateNewSong(Write http.ResponseWriter, Req *http.Request) {
	var SongRequest request.SongRequest
	songDataStr := Req.FormValue("songData")
	fmt.Println("Received songData:", songDataStr) // 🔍 In ra để debug

	errorToConvert := json.Unmarshal([]byte(Req.FormValue("songData")), &SongRequest)
	fmt.Println(SongRequest)
	Files, _, errorToGetFile := Req.FormFile("file")
	SongFile := request.SongFile{
		File: Files,
	}
	if errorToGetFile != nil {

		http.Error(Write, "failed to File", http.StatusBadRequest)
		return
	}
	if errorToConvert != nil {
		fmt.Print(errorToConvert)
		http.Error(Write, "failed to Json", http.StatusBadRequest)
		return
	}
	resp, errorToCreateSong := Controller.songService.CreateNewSong(SongRequest, SongFile)
	if errorToCreateSong != nil {
		http.Error(Write, "failed to Song", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(resp)

}
