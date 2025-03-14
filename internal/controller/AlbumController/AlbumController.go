package albumcontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"ten_module/internal/DTO/request"
	"ten_module/internal/service/albumservice"

	"github.com/gorilla/mux"
)

type AlbumController struct {
	AlbumServe *albumservice.AlbumSerivce
}

var AlbumControll *AlbumController

func InitAlbumController() {
	AlbumControll = &AlbumController{
		AlbumServe: albumservice.AlbumServe,
	}
}
func (AlbumControll *AlbumController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/createalbum", AlbumControll.CreateAlbum).Methods("POST")

}
func (AlbumControll *AlbumController) CreateAlbum(Write http.ResponseWriter, Request *http.Request) {
	var AlbumReq request.AlbumRequest
	var SongFileReq []request.SongFileAlbum
	AlbumRequest := Request.FormValue("album_request")
	ErrorToConvert := json.Unmarshal([]byte(AlbumRequest), &AlbumReq)
	if ErrorToConvert != nil {
		log.Print("ss")
		log.Print(ErrorToConvert)
		http.Error(Write, "failed to convert to json", http.StatusBadRequest)
		return
	}
	SongFile := Request.MultipartForm.File["song_file"]
	log.Print(SongFile)
	SongReq := AlbumReq.Song
	if len(SongFile) == 0 {
		log.Print("Lỗi: Không có tệp nào được tải lên!")
		http.Error(Write, "Không có tệp bài hát nào được tải lên", http.StatusBadRequest)
		return
	}

	if len(SongFile) != len(SongReq) {
		log.Print("Lỗi: Số lượng file tải lên không khớp với số lượng bài hát!")
		http.Error(Write, "Số lượng file không khớp với số lượng bài hát", http.StatusBadRequest)
		return
	}
	for i := 0; i < len(SongReq); i++ {
		file, err := SongFile[i].Open()
		if err != nil {
			http.Error(Write, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		SongFileReq = append(SongFileReq, request.SongFileAlbum{
			SongName: SongReq[i].NameSong,
			File:     file,
		})
	}
	Resp, ErrorToCreateAlbum := AlbumControll.AlbumServe.CreateAlbum(AlbumReq, SongFileReq)
	if ErrorToCreateAlbum != nil {
		log.Print(ErrorToCreateAlbum)
		http.Error(Write, "ErrorToCreateAlbum", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(Write).Encode(Resp)
}
