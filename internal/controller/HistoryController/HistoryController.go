package historycontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	middleware "ten_module/Middleware"
	"ten_module/internal/service/listenhistoryservice"

	"github.com/gorilla/mux"
)

type HistoryController struct {
	HistoryService *listenhistoryservice.ListenHistoryService
	middleware     *middleware.UseMiddleware
}

var HistoryControllers *HistoryController

func InitHistoryControllers() {
	HistoryControllers = &HistoryController{
		middleware:     middleware.Middlewares,
		HistoryService: listenhistoryservice.HistoryService,
	}
}
func (HistoryControll *HistoryController) RegisterRoute(r *mux.Router) {
	r.HandleFunc("/savehistory", HistoryControll.SaveHistoryListen).Methods("POST")
}
func (HistoryControll *HistoryController) SaveHistoryListen(Write http.ResponseWriter, Req *http.Request) {
	UserId := Req.URL.Query().Get("userid")
	SongId := Req.URL.Query().Get("songid")

	SongIdConvert, ErrorToConvertString := strconv.Atoi(SongId)
	if ErrorToConvertString != nil {
		http.Error(Write, "failed to Convert", http.StatusBadRequest)
		log.Print(ErrorToConvertString)
		return
	}
	Resp, ErrorToSave := HistoryControll.HistoryService.SaveHistoryListen(UserId, SongIdConvert)
	if ErrorToSave != nil {
		http.Error(Write, "failed to Convert", http.StatusBadRequest)
		log.Print(ErrorToSave)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
