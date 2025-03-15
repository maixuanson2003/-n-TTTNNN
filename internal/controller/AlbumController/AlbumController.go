package albumcontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	r.HandleFunc("/album/getlist", AlbumControll.GetListAlbum).Methods("GET")
	r.HandleFunc("/album/{id}", AlbumControll.GetAlbumById).Methods("GET")
	r.HandleFunc("/getalbum/artist", AlbumControll.GetAlbumByArtist).Methods("GET")

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
func (AlbumControll *AlbumController) GetListAlbum(Write http.ResponseWriter, Request *http.Request) {
	Resp, ErrorToGetAlbum := AlbumControll.AlbumServe.GetListAlbum()
	if ErrorToGetAlbum != nil {
		log.Print(ErrorToGetAlbum)
		http.Error(Write, "faile to get list", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(Write).Encode(Resp)
}
func (AlbumControll *AlbumController) GetAlbumById(Write http.ResponseWriter, Request *http.Request) {
	url := Request.URL.Path
	TakeAlbumId := strings.Split(url, "/")[3]
	AlbumId, ErrorToConvertToNumber := strconv.Atoi(TakeAlbumId)
	if ErrorToConvertToNumber != nil {
		log.Print(ErrorToConvertToNumber)
		http.Error(Write, "failed to convert to int", http.StatusBadRequest)
		return

	}
	Resp, ErrorToGetAlbum := AlbumControll.AlbumServe.GetAlbumById(AlbumId)
	if ErrorToGetAlbum != nil {
		log.Print(ErrorToGetAlbum)
		http.Error(Write, "failed to get album", http.StatusBadRequest)
		return

	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(Write).Encode(Resp)
}
func (AlbumControll *AlbumController) GetAlbumByArtist(Write http.ResponseWriter, Request *http.Request) {
	ArtistId, ErrorToConvert := strconv.Atoi(Request.URL.Query().Get("artistid"))
	if ErrorToConvert != nil {
		log.Print(ErrorToConvert)
		http.Error(Write, "failed to convert", http.StatusBadRequest)
		return

	}

	Resp, ErrorToGetAlbum := AlbumControll.AlbumServe.GetAlbumByArtist(ArtistId)
	if ErrorToGetAlbum != nil {
		log.Print(ErrorToGetAlbum)
		http.Error(Write, "failed get album by  artist", http.StatusBadRequest)
		return

	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(Write).Encode(Resp)

}
