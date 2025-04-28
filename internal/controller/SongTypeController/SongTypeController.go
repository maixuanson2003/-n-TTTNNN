package songtypecontroller

import (
	"encoding/json"
	"net/http"
	"ten_module/internal/service/songtypeservice"

	"github.com/gorilla/mux"
)

type SongTypeController struct {
	SongTypeServe *songtypeservice.SongTypeService
}

var SongTypeControll *SongTypeController

func InitSongTypeController() {
	SongTypeControll = &SongTypeController{
		SongTypeServe: songtypeservice.SongTypeServe,
	}
}
func (SongTypeControll *SongTypeController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/listtype", SongTypeControll.GetListType).Methods("GET")
	r.HandleFunc("/createtype", SongTypeControll.CreateType).Methods("POST")

}
func (SongTypeControll *SongTypeController) GetListType(Write http.ResponseWriter, Request *http.Request) {
	Resp, err := SongTypeControll.SongTypeServe.GetListSongType()
	if err != nil {
		http.Error(Write, "Invalid request payload", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
func (SongTypeControll *SongTypeController) CreateType(Write http.ResponseWriter, Request *http.Request) {
	SongType := Request.URL.Query().Get("type")
	Resp, err := SongTypeControll.SongTypeServe.CreateSongType(SongType)
	if err != nil {
		http.Error(Write, "failed to create Song Type", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)

}
