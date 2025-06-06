package songtypecontroller

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	r.HandleFunc("/updatetype", SongTypeControll.UpdateType).Methods("PUT")
	r.HandleFunc("/deletetype/{id}", SongTypeControll.DeleteType).Methods("DELETE")
	r.HandleFunc("/gettype/{id}", SongTypeControll.GetType).Methods("GET")

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
func (SongTypeControll *SongTypeController) GetType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	resp, err := SongTypeControll.SongTypeServe.GetTypeById(id)
	if err != nil {
		http.Error(w, "Failed to get song type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
func (SongTypeControll *SongTypeController) UpdateType(w http.ResponseWriter, r *http.Request) {
	typeParam := r.URL.Query().Get("type")
	idParam := r.URL.Query().Get("id")

	if typeParam == "" || idParam == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	// Convert id string to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	resp, err := SongTypeControll.SongTypeServe.UpdateSongType(typeParam, id)
	if err != nil {
		http.Error(w, "Failed to update song type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Xoá loại bài hát theo ID
func (SongTypeControll *SongTypeController) DeleteType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	resp := SongTypeControll.SongTypeServe.DeleteTypeById(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
