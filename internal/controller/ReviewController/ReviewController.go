package reviewcontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	middleware "ten_module/Middleware"
	"ten_module/internal/DTO/request"
	"ten_module/internal/service/reviewservice"

	"github.com/gorilla/mux"
)

type ReviewController struct {
	ReviewServe *reviewservice.ReviewService
	MiddleWare  *middleware.UseMiddleware
}

var ReviewControll *ReviewController

func InitReviewController() {
	ReviewControll = &ReviewController{
		ReviewServe: reviewservice.ReviewServe,
		MiddleWare:  middleware.Middlewares,
	}
}
func (ReviewControll *ReviewController) RegisterRoute(r *mux.Router) {
	middleware := ReviewControll.MiddleWare
	r.HandleFunc("/reviewlist", ReviewControll.GetListReview).Methods("GET")
	r.HandleFunc("/reviewlistinsong/{id}", ReviewControll.GetListReviewBySong).Methods("GET")
	r.HandleFunc("/review/user", ReviewControll.GetListReviewForUser).Methods("GET")
	r.HandleFunc("/createreview", middleware.Chain(ReviewControll.CreateReview, middleware.CheckToken(), middleware.VerifyRole([]string{"USER"}))).Methods("POST")
	r.HandleFunc("/updatereview", middleware.Chain(ReviewControll.UpdateReview, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("PUT")
	r.HandleFunc("/deletereview/{id}", middleware.Chain(ReviewControll.DeleteReview, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN", "USER"}))).Methods("DELETE")
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
	log.Print("ssss")
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
func (reviewController *ReviewController) UpdateReview(w http.ResponseWriter, r *http.Request) {

	reviewIdStr := r.URL.Query().Get("reviewid")
	reviewId, err := strconv.Atoi(reviewIdStr)
	status := r.URL.Query().Get("status")
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	response, err := reviewController.ReviewServe.UpdateReview(status, reviewId)
	if err != nil {
		http.Error(w, response.Message, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func (reviewController *ReviewController) DeleteReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewIdStr := vars["id"]
	reviewId, err := strconv.Atoi(reviewIdStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	response, err := reviewController.ReviewServe.DeleteReviewById(reviewId)
	if err != nil {
		http.Error(w, response.Message, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func (reviewController *ReviewController) DeleteReviewForUser(w http.ResponseWriter, r *http.Request) {

	reviewIdStr := r.URL.Query().Get("reviewid")
	reviewId, err := strconv.Atoi(reviewIdStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}
	userId := r.URL.Query().Get("userid")
	if userId == "" {
		http.Error(w, "Unauthorized: missing user ID", http.StatusUnauthorized)
		return
	}
	response, err := reviewController.ReviewServe.DeleteReviewByIdForUser(reviewId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func (reviewController *ReviewController) GetListReviewForUser(w http.ResponseWriter, r *http.Request) {

	userId := r.URL.Query().Get("userid")
	if userId == "" {
		http.Error(w, "Unauthorized: missing user ID", http.StatusUnauthorized)
		return
	}
	response, err := reviewController.ReviewServe.GetListReviewForUser(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
