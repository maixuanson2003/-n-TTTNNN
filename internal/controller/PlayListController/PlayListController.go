package playlistcontroller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"ten_module/internal/service/playlistservice"

	"github.com/gorilla/mux"
)

type PlayListController struct {
	PlayListServ *playlistservice.PlayListService
}

var PlayListControll *PlayListController

func InitPlayListController() {
	PlayListControll = &PlayListController{
		PlayListServ: playlistservice.PlayListServe,
	}
}
func (Controller *PlayListController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/playlist/{id}", Controller.GetPlayListByUser).Methods("GET")
	r.HandleFunc("/songplay/{id}", Controller.GetSongByPlayList).Methods("GET")
	r.HandleFunc("/createplay", Controller.CreatePlayList).Methods("POST")
	r.HandleFunc("/addsong", Controller.AddSongToPlayList).Methods("PUT")

}
func (Controller *PlayListController) GetPlayListByUser(Write http.ResponseWriter, Request *http.Request) {
	url := Request.URL.Path
	userId := strings.Split(url, "/")[3]
	PlayList, ErrorToGetPlayList := Controller.PlayListServ.GetPlayListByUser(userId)
	if ErrorToGetPlayList != nil {
		http.Error(Write, "not found", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(PlayList)
}
func (Controller *PlayListController) GetSongByPlayList(Write http.ResponseWriter, Request *http.Request) {
	url := Request.URL.Path
	Id := strings.Split(url, "/")[3]
	PlayListId, errorToConvert := strconv.Atoi(Id)
	if errorToConvert != nil {
		http.Error(Write, "Failed to convert", http.StatusBadRequest)
		return
	}
	Resp, ErrorToGetSong := Controller.PlayListServ.GetSongByPLayList(PlayListId)
	if ErrorToGetSong != nil {
		http.Error(Write, "Failed to convert", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
func (Controller *PlayListController) CreatePlayList(Write http.ResponseWriter, Request *http.Request) {
	UserId := Request.URL.Query().Get("userid")
	NamePlayList := Request.URL.Query().Get("nameplaylist")

	Resp, ErrorToCreatePlaylist := Controller.PlayListServ.CreatePlayList(NamePlayList, UserId)
	if ErrorToCreatePlaylist != nil {
		http.Error(Write, "Failed to convert", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
func (Controller *PlayListController) AddSongToPlayList(Write http.ResponseWriter, Request *http.Request) {
	SongId, ErrorToConvertSongId := strconv.Atoi(Request.URL.Query().Get("songid"))
	if ErrorToConvertSongId != nil {
		http.Error(Write, "Failed to convert", http.StatusBadRequest)
		return
	}
	PlayListId, ErrorToConverPlayListId := strconv.Atoi(Request.URL.Query().Get("playlistid"))
	if ErrorToConverPlayListId != nil {
		http.Error(Write, "Failed to convert", http.StatusBadRequest)
		return
	}
	Resp, ErrorToAddSong := Controller.PlayListServ.AddSongToPlayList(SongId, PlayListId)
	if ErrorToAddSong != nil {
		http.Error(Write, "Failed to add song", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
