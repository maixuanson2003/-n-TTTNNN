package songcontroller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	r.HandleFunc("/song/{id}", Controller.DownLoadSong).Methods("GET")
	r.HandleFunc("/Like", Controller.UserLikeSong).Methods("POST")
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
func (Controller *SongController) DownLoadSong(Write http.ResponseWriter, Req *http.Request) {
	url := Req.URL.Path
	fmt.Print("ssss")
	GetSongId := strings.Split(url, "/")[3]
	SongId, ErrorToConvertString := strconv.Atoi(GetSongId)
	if ErrorToConvertString != nil {
		http.Error(Write, "failed to Convert", http.StatusBadRequest)
		log.Print(ErrorToConvertString)
		return
	}
	resp, errorToHandleDownLoad := Controller.songService.DownLoadSong(SongId)
	fmt.Print("ssss")
	if errorToHandleDownLoad != nil {
		http.Error(Write, "failed to get download", http.StatusBadRequest)
		log.Print(ErrorToConvertString)
		return
	}
	defer resp.Resp.Body.Close()
	Write.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", resp.NameSong))
	Write.Header().Set("Content-Type", "audio/mpeg")
	Write.Header().Set("Content-Transfer-Encoding", "binary")
	_, errorToConvert := io.Copy(Write, resp.Resp.Body)
	if errorToConvert != nil {
		log.Print(errorToConvert)
		return
	}

}
func (Controller *SongController) UserLikeSong(Write http.ResponseWriter, Req *http.Request) {
	UserId := Req.URL.Query().Get("userid")
	SongId := Req.URL.Query().Get("songid")
	fmt.Print("sss")
	SongIdConvert, ErrorToConvertString := strconv.Atoi(SongId)
	if ErrorToConvertString != nil {
		http.Error(Write, "failed to Convert", http.StatusBadRequest)
		log.Print(ErrorToConvertString)
		return
	}
	resp, ErrorToLike := Controller.songService.UserLikeSong(SongIdConvert, UserId)
	if ErrorToLike != nil {
		http.Error(Write, "failed to Convert", http.StatusBadRequest)
		log.Print(ErrorToLike)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(resp)

}
