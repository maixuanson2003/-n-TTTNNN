package songcontroller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	r.HandleFunc("/foruser/{id}", Controller.GetSongForUser).Methods("GET")
	r.HandleFunc("/geturl", Controller.GetAllUrlSong).Methods("GET")
	r.HandleFunc("/getsongall", Controller.GetAllSong).Methods("GET")
	r.HandleFunc("/recommend", Controller.GetSongForUserRecommend).Methods("GET")
	r.HandleFunc("/search", Controller.SearchSongByKeyWord).Methods("GET")
	r.HandleFunc("/filtersong", Controller.FilterSong).Methods("GET")
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
func (Controller *SongController) GetSongForUser(Write http.ResponseWriter, Req *http.Request) {
	url := Req.URL.Path
	UserId := strings.Split(url, "/")[3]
	Resp, ErrorToGetSong := Controller.songService.GetListSongForUser(UserId)
	if ErrorToGetSong != nil {
		http.Error(Write, "failed to get song", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
func (Controller *SongController) GetAllUrlSong(Write http.ResponseWriter, Req *http.Request) {
	FolderPath := "C:\\Users\\DPC\\Desktop\\MusicMp4\\internal\\music"
	FileArray, ErrorToGetFile := os.ReadDir(FolderPath)
	if ErrorToGetFile != nil {
		http.Error(Write, "write file false", http.StatusBadRequest)
		return
	}
	nameFile := []string{}
	for _, itemFile := range FileArray {
		FileName := itemFile.Name()
		url := "http://localhost:8080/music/" + FileName
		nameFile = append(nameFile, url)

	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(nameFile)

}
func (Controller *SongController) GetAllSong(Write http.ResponseWriter, Req *http.Request) {
	page := Req.URL.Query().Get("page")
	finalPage, errorToConvert := strconv.Atoi(page)
	if errorToConvert != nil {
		http.Error(Write, "faile to convert", http.StatusBadRequest)
		return
	}
	SongResponse, errorToGetListSong := Controller.songService.GetAllSong(finalPage)
	if errorToGetListSong != nil {
		http.Error(Write, "write file false", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(SongResponse)

}
func (Controller *SongController) GetSongForUserRecommend(Write http.ResponseWriter, Req *http.Request) {
	UserId := Req.URL.Query().Get("userid")
	Resp, ErrorToGetSong := Controller.songService.GetSongForUser(UserId)
	if ErrorToGetSong != nil {
		http.Error(Write, "faile", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
func (Controller *SongController) SearchSongByKeyWord(Write http.ResponseWriter, Req *http.Request) {
	KeyWord := Req.URL.Query().Get("keyword")
	Resp, ErrorToSearchSong := Controller.songService.SearchSongByKeyWord(KeyWord)
	if ErrorToSearchSong != nil {
		http.Error(Write, "faile to search song", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
func (Controller *SongController) FilterSong(Write http.ResponseWriter, Req *http.Request) {
	artistIdsStr := Req.URL.Query().Get("artistId")
	typeIdsStr := Req.URL.Query().Get("typeId")

	// Xử lý artistIds
	var artistIds []int
	if artistIdsStr != "" {
		artistIdsStrArr := strings.Split(artistIdsStr, ",")
		for _, idStr := range artistIdsStrArr {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(Write, "Invalid artistId", http.StatusBadRequest)
				return
			}
			artistIds = append(artistIds, id)
		}
	}

	// Xử lý typeIds
	var typeIds []int
	if typeIdsStr != "" {
		typeIdsStrArr := strings.Split(typeIdsStr, ",")
		for _, idStr := range typeIdsStrArr {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(Write, "Invalid typeId", http.StatusBadRequest)
				return
			}
			typeIds = append(typeIds, id)
		}
	}

	songs, err := Controller.songService.FilterSong(artistIds, typeIds)
	if err != nil {
		http.Error(Write, fmt.Sprintf("Error filtering songs: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(songs)
}
