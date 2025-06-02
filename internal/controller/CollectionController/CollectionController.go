package collectioncontroller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	middleware "ten_module/Middleware"
	collectionservice "ten_module/internal/service/CollectionService"

	"github.com/gorilla/mux"
)

type CollectionController struct {
	CollectionService *collectionservice.CollectionService
	MiddleWare        *middleware.UseMiddleware
}

var CollectionControll *CollectionController

func InitCollectionController() {
	CollectionControll = &CollectionController{
		CollectionService: collectionservice.CollectionServe,
		MiddleWare:        middleware.Middlewares,
	}
}
func (CollectionControll *CollectionController) RegisterRoute(r *mux.Router) {
	middleware := CollectionControll.MiddleWare
	r.HandleFunc("/listcollect", CollectionControll.GetListCollection).Methods("GET")
	r.HandleFunc("/collect/{id}", CollectionControll.GetCollectionById).Methods("GET")
	r.HandleFunc("/createcollect", middleware.Chain(CollectionControll.CreateCollection, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("POST")
	r.HandleFunc("/addtocollect", middleware.Chain(CollectionControll.AddSongToCollection, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("PUT")
	r.HandleFunc("/deletesongcollect", middleware.Chain(CollectionControll.DeleteSongFromCollect, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("DELETE")
	r.HandleFunc("/deletecollect/{id}", middleware.Chain(CollectionControll.DeleteCollection, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("DELETE")
}

func (CollectionControll *CollectionController) GetListCollection(Write http.ResponseWriter, Request *http.Request) {
	Resp, ErrorToGetCollect := CollectionControll.CollectionService.GetListCollection()
	if ErrorToGetCollect != nil {
		http.Error(Write, "failed to get List collect", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(Write).Encode(Resp)
}
func (CollectionControll *CollectionController) GetCollectionById(Write http.ResponseWriter, Request *http.Request) {
	url := Request.URL.Path
	GetCollectionId := strings.Split(url, "/")[3]
	CollectionId, ErrorToConvertString := strconv.Atoi(GetCollectionId)
	if ErrorToConvertString != nil {
		http.Error(Write, "failed to get  collect", http.StatusBadRequest)
		return
	}
	Resp, ErrorToGetCollect := CollectionControll.CollectionService.GetCollectionById(CollectionId)
	if ErrorToGetCollect != nil {
		http.Error(Write, "failed to get  collect", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(Write).Encode(Resp)
}
func (CollectionControll *CollectionController) CreateCollection(Write http.ResponseWriter, Request *http.Request) {
	NameCollection := Request.URL.Query().Get("namecollection")
	if NameCollection == "" {
		http.Error(Write, "require name collect", http.StatusBadRequest)
		return
	}
	Resp, ErrorToCreateCollect := CollectionControll.CollectionService.CreateCollection(NameCollection)
	if ErrorToCreateCollect != nil {
		http.Error(Write, "Failed to create collect", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(Write).Encode(Resp)
}
func (CollectionControll *CollectionController) AddSongToCollection(Write http.ResponseWriter, Request *http.Request) {
	SongId, ErrorToConvertSongId := strconv.Atoi(Request.URL.Query().Get("songid"))
	if ErrorToConvertSongId != nil {
		http.Error(Write, "Failed to convert", http.StatusBadRequest)
		return
	}
	CollectionId, ErrorToConverCollectionId := strconv.Atoi(Request.URL.Query().Get("collectionid"))
	if ErrorToConverCollectionId != nil {
		http.Error(Write, "Failed to convert", http.StatusBadRequest)
		return
	}
	Resp, ErrorToAddSongToCollect := CollectionControll.CollectionService.AddSongToCollection(SongId, CollectionId)
	if ErrorToAddSongToCollect != nil {
		http.Error(Write, "Failed to add song to collection", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(Write).Encode(Resp)

}
func (CollectionControll *CollectionController) DeleteSongFromCollect(Write http.ResponseWriter, Request *http.Request) {
	SongId, ErrorToConvertSongId := strconv.Atoi(Request.URL.Query().Get("songid"))
	if ErrorToConvertSongId != nil {
		http.Error(Write, "Failed to convert", http.StatusBadRequest)
		return
	}
	CollectionId, ErrorToConverCollectionId := strconv.Atoi(Request.URL.Query().Get("collectionid"))
	if ErrorToConverCollectionId != nil {
		http.Error(Write, "Failed to convert", http.StatusBadRequest)
		return
	}
	Resp, ErrorToAddSongToCollect := CollectionControll.CollectionService.DeleteSongFromCollection(SongId, CollectionId)
	if ErrorToAddSongToCollect != nil {
		http.Error(Write, "Failed to add song to collection", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(Write).Encode(Resp)
}
func (CollectionControll *CollectionController) DeleteCollection(Write http.ResponseWriter, Request *http.Request) {
	url := Request.URL.Path
	GetCollectionId := strings.Split(url, "/")[3]
	CollectionId, ErrorToConvertString := strconv.Atoi(GetCollectionId)
	if ErrorToConvertString != nil {
		http.Error(Write, "failed to get  collect", http.StatusBadRequest)
		return
	}
	resp, err := CollectionControll.CollectionService.DeleteCollectionByIds(CollectionId)
	if err != nil {
		http.Error(Write, "Failed to delete collection", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusAccepted)
	json.NewEncoder(Write).Encode(resp)
}
