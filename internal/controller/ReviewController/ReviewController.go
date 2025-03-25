package reviewcontroller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"ten_module/internal/DTO/request"
	"ten_module/internal/service/reviewservice"

	"github.com/gorilla/mux"
)

type ReviewController struct {
	ReviewServe *reviewservice.ReviewService
}

var ReviewControll *ReviewController

func InitReviewController() {
	ReviewControll = &ReviewController{
		ReviewServe: reviewservice.ReviewServe,
	}
}
func (ReviewControll *ReviewController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/reviewlist", ReviewControll.GetListReview).Methods("GET")
	r.HandleFunc("/reviewlistinsong/{id}", ReviewControll.GetListReviewBySong).Methods("GET")
	r.HandleFunc("/createreview", ReviewControll.CreateReview).Methods("POST")

}
func (ReviewControll *ReviewController) CreateReview(Write http.ResponseWriter, Request *http.Request) {
	var ReviewRequest *request.ReviewRequest
	json.NewDecoder(Request.Body).Decode(&ReviewRequest)
	Resp, ErrorToCreateReview := ReviewControll.ReviewServe.CreateReview(*ReviewRequest)
	if ErrorToCreateReview != nil {
		http.Error(Write, "failed to create review", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)

}
func (ReviewControll *ReviewController) GetListReview(Write http.ResponseWriter, Request *http.Request) {
	Resp, ErrorToGetListReview := ReviewControll.ReviewServe.GetListReview()
	if ErrorToGetListReview != nil {
		http.Error(Write, "faile to get list review", http.StatusBadRequest)
		return

	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
func (ReviewControll *ReviewController) GetListReviewBySong(Write http.ResponseWriter, Request *http.Request) {
	url := Request.URL.Path
	Id := strings.Split(url, "/")[3]
	SongId, errorToConvert := strconv.Atoi(Id)
	if errorToConvert != nil {
		http.Error(Write, "failed to convert", http.StatusBadRequest)
		return
	}
	Resp, ErrorToGetListReview := ReviewControll.ReviewServe.GetListReviewBySong(SongId)
	if ErrorToGetListReview != nil {
		http.Error(Write, "failed to convert", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
