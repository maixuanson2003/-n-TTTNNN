package songcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	middleware "ten_module/Middleware"
	"ten_module/internal/DTO/request"
	openai "ten_module/internal/Helper/openAi"
	"ten_module/internal/service/songservice"
	"time"

	"github.com/go-playground/validator/v10"
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
	middleware := Controller.middlware
	r.HandleFunc("/song/create", middleware.Chain(Controller.CreateNewSong, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("POST")
	r.HandleFunc("/download/song", Controller.DownloadHandler).Methods("GET")
	r.HandleFunc("/getsong/{id}", Controller.GetSongById).Methods("GET")
	r.HandleFunc("/updatesong/{id}", middleware.Chain(Controller.UpdateSong, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("PUT")
	r.HandleFunc("/Like", Controller.UserLikeSong).Methods("POST")
	r.HandleFunc("/geturl", Controller.GetAllUrlSong).Methods("GET")
	r.HandleFunc("/getsongall", Controller.GetAllSong).Methods("GET")
	r.HandleFunc("/recommend", Controller.GetSongForUserRecommend).Methods("GET")
	r.HandleFunc("/search", Controller.SearchSongByKeyWord).Methods("GET")
	r.HandleFunc("/filtersong", Controller.FilterSong).Methods("GET")
	r.HandleFunc("/delete/song", middleware.Chain(Controller.DeleteSongById, middleware.CheckToken(), middleware.VerifyRole([]string{"ADMIN"}))).Methods("DELETE")
	r.HandleFunc("/topweek/song", Controller.GetTopSongsThisWeek).Methods("GET")
	// r.HandleFunc("/update/song/album", Controller.UpdateSongAlbum).Methods("POST")
	r.HandleFunc("/chatbot/song", Controller.ChatMusicHandler).Methods("POST")
	r.HandleFunc("/song/week/chart", Controller.GetWeeklyChartDataPerDay).Methods("GET")
	r.HandleFunc("/compe/song/top", Controller.GetTopSongCompe).Methods("GET")
	r.HandleFunc("/getrecommed/song", Controller.GetSimilarSongsRecommend).Methods("GET")
	r.HandleFunc("/get/like/{id}", Controller.GetSongLikeByUser).Methods("GET")
	r.HandleFunc("/dishlike", Controller.UserDishLikeSong).Methods("POST")
	r.HandleFunc("/list/song", Controller.GetSongList).Methods("GET")
}

var validate = validator.New()

func (Controller *SongController) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("url")
	if origin == "" {
		http.Error(w, "missing url query parameter", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, origin, nil)
	if err != nil {
		http.Error(w, "bad origin URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "cannot reach origin: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "origin returned "+resp.Status, http.StatusBadGateway)
		return
	}
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		w.Header().Set("Content-Disposition", cd)
	} else {
		w.Header().Set("Content-Disposition", "attachment")
	}
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, resp.Body)
}
func (Controller *SongController) GetTopSongCompe(w http.ResponseWriter, r *http.Request) {
	ranges := r.URL.Query().Get("range")
	log.Print(ranges)
	topSongs := Controller.songService.GetBookTopRange(ranges)

	w.Header().Set("Content-Type", "application/json")
	if topSongs == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Không thể lấy danh sách bài hát",
		})
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    topSongs,
	})
}
func (Controller *SongController) GetWeeklyChartDataPerDay(w http.ResponseWriter, r *http.Request) {
	topN := 3
	query := r.URL.Query()
	if val := query.Get("top"); val != "" {
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			topN = n
		}
	}
	chartData := Controller.songService.GetWeeklyChartDataPerDay(topN)
	w.Header().Set("Content-Type", "application/json")
	if chartData == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Không thể lấy dữ liệu chart",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    chartData,
	})
}
func (Controller *SongController) ChatMusicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
		return
	}
	message := r.URL.Query().Get("message")

	query, err := openai.ExtractMusicInfo(message)
	if err != nil {
		log.Print(err)
		http.Error(w, fmt.Sprintf("Lỗi GPT: %v", err), http.StatusInternalServerError)
		return
	}
	song := Controller.songService.GetDataFromPromb(query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// func (Controller *SongController) UpdateSongAlbum(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseMultipartForm(32 << 20)
// 	if err != nil {
// 		http.Error(w, "Không thể đọc multipart form", http.StatusBadRequest)
// 		return
// 	}
// 	albumIdStr := r.URL.Query().Get("albumid")
// 	if albumIdStr == "" {
// 		http.Error(w, "albumId là bắt buộc", http.StatusBadRequest)
// 		return
// 	}
// 	var albumId int
// 	fmt.Sscanf(albumIdStr, "%d", &albumId)
// 	songsJson := r.FormValue("songs")
// 	if songsJson == "" {
// 		log.Print("looi")
// 		http.Error(w, "Thiếu thông tin bài hát (songs)", http.StatusBadRequest)
// 		return
// 	}

// 	var songRequests []request.SongRequest
// 	err = json.Unmarshal([]byte(songsJson), &songRequests)
// 	if err != nil {
// 		log.Print(err)
// 		http.Error(w, "Không thể parse JSON bài hát", http.StatusBadRequest)
// 		return
// 	}
// 	SongFile := r.MultipartForm.File["file"]
// 	if len(SongFile) != len(songRequests) {
// 		log.Print("do dai ko khop")
// 		http.Error(w, "Không thể parse JSON bài hát", http.StatusBadRequest)
// 		return
// 	}
// 	songFiles := []request.SongFile{}

// 	for i := 0; i < len(SongFile); i++ {
// 		file, err := SongFile[i].Open()
// 		if err != nil {
// 			log.Print(err)
// 			http.Error(w, "Failed to open file", http.StatusInternalServerError)
// 			return
// 		}
// 		defer file.Close()
// 		songFiles = append(songFiles, request.SongFile{File: file})
// 	}

// 	// Gọi service để xử lý
// 	Controller.songService.UpdateSongAlbum(&albumId, songRequests, songFiles)
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": "Đang cập nhật danh sách bài hát cho album...",
// 	})
// }

func (Controller *SongController) GetTopSongsThisWeek(w http.ResponseWriter, r *http.Request) {
	topSongs := Controller.songService.GetBookTopRange("week")

	w.Header().Set("Content-Type", "application/json")
	if topSongs == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Không thể lấy danh sách bài hát",
		})
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    topSongs,
	})
}
func (Controller *SongController) CreateNewSong(Write http.ResponseWriter, Req *http.Request) {
	var SongRequest request.SongRequest
	songDataStr := Req.FormValue("songData")
	fmt.Println("Received songData:", songDataStr)

	errorToConvert := json.Unmarshal([]byte(Req.FormValue("songData")), &SongRequest)
	errorsToValidate := validate.Struct(SongRequest)
	if errorsToValidate != nil {
		validationErrors := errorsToValidate.(validator.ValidationErrors)
		var errorMsg string
		for _, e := range validationErrors {
			errorMsg += fmt.Sprintf("Trường '%s' không hợp lệ (%s); ", e.Field(), e.Tag())
		}
		log.Print(errorMsg)
		http.Error(Write, errorMsg, http.StatusBadRequest)

		return
	}
	fmt.Println(SongRequest)
	_, Header, errorToGetFile := Req.FormFile("file")
	SongFile := request.SongFile{
		File: Header,
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
func (Controller *SongController) UpdateSong(Write http.ResponseWriter, Req *http.Request) {
	var SongRequest request.SongRequest
	songDataStr := Req.FormValue("songData")
	fmt.Println("Received songData:", songDataStr)
	url := Req.URL.Path
	GetSongId := strings.Split(url, "/")[3]
	SongId, ErrorToConvertString := strconv.Atoi(GetSongId)
	if ErrorToConvertString != nil {
		http.Error(Write, "failed to Convert", http.StatusBadRequest)
		log.Print(ErrorToConvertString)
		return
	}
	errorToConvert := json.Unmarshal([]byte(Req.FormValue("songData")), &SongRequest)
	if errorToConvert != nil {
		log.Print(errorToConvert)
		http.Error(Write, "failed to Json", http.StatusBadRequest)
		return
	}
	var SongFile request.SongFile
	_, fileHeader, fileErr := Req.FormFile("file")
	SongFile.File = fileHeader
	if fileErr != nil && !errors.Is(fileErr, http.ErrMissingFile) {
		log.Print(fileErr)
		http.Error(Write, "failed to get file", http.StatusBadRequest)
		return
	}
	resp, errToUpdate := Controller.songService.UpdateSong(SongRequest, SongId, SongFile)
	if errToUpdate != nil {
		log.Print(errToUpdate)
		http.Error(Write, "failed to update", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(resp)

}
func (Controller *SongController) GetSongList(Write http.ResponseWriter, Req *http.Request) {

	resp, errorToGetSong := Controller.songService.GetListSong()
	if errorToGetSong != nil {
		http.Error(Write, "failed to get Song", http.StatusBadRequest)
		log.Print(errorToGetSong)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(resp)

}
func (Controller *SongController) GetSongById(Write http.ResponseWriter, Req *http.Request) {
	url := Req.URL.Path
	fmt.Print("ssss")
	GetSongId := strings.Split(url, "/")[3]
	SongId, ErrorToConvertString := strconv.Atoi(GetSongId)
	if ErrorToConvertString != nil {
		http.Error(Write, "failed to Convert", http.StatusBadRequest)
		log.Print(ErrorToConvertString)
		return
	}
	resp, errorToGetSong := Controller.songService.GetSongById(SongId)
	if errorToGetSong != nil {
		http.Error(Write, "failed to get Song", http.StatusBadRequest)
		log.Print(errorToGetSong)
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
func (c *SongController) GetSongLikeByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"] // {userId} trong route
	if userID == "" {    // không tìm thấy hoặc rỗng
		http.Error(w, "missing userId in path", http.StatusBadRequest)
		return
	}

	songs, err := c.songService.GetSongByUserId(userID)
	if err != nil {
		http.Error(w, "failed to get songs for user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(songs)
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
func (Controller *SongController) UserDishLikeSong(Write http.ResponseWriter, Req *http.Request) {
	UserId := Req.URL.Query().Get("userid")
	SongId := Req.URL.Query().Get("songid")
	fmt.Print("sss")
	SongIdConvert, ErrorToConvertString := strconv.Atoi(SongId)
	if ErrorToConvertString != nil {
		http.Error(Write, "failed to Convert", http.StatusBadRequest)
		log.Print(ErrorToConvertString)
		return
	}
	resp, ErrorToLike := Controller.songService.UserDislikeSong(SongIdConvert, UserId)
	if ErrorToLike != nil {
		http.Error(Write, "failed to Convert", http.StatusBadRequest)
		log.Print(ErrorToLike)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(resp)

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

	SongResponse, errorToGetListSong := Controller.songService.GetAllSong()
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
	// Resp, ErrorToGetSong := Controller.songService.GetSongForUser(UserId)
	Resp, ErrorToGetSong := Controller.songService.GetSongForUserV2(UserId, 7, 7)
	if ErrorToGetSong != nil {
		http.Error(Write, "faile", http.StatusBadRequest)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(Resp)
}
func (Controller *SongController) GetSimilarSongsRecommend(w http.ResponseWriter, r *http.Request) {
	songIdStr := r.URL.Query().Get("songid")
	if songIdStr == "" {
		http.Error(w, "missing songid", http.StatusBadRequest)
		return
	}
	songId, err := strconv.Atoi(songIdStr)
	if err != nil {
		http.Error(w, "invalid songid", http.StatusBadRequest)
		return
	}
	resp, err := Controller.songService.GetSimilarSongs(songId)
	if err != nil {
		http.Error(w, "failed to get recommendations", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
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
	print(artistIdsStr)

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
func (Controller *SongController) DeleteSongById(Write http.ResponseWriter, Req *http.Request) {
	songidparam := Req.URL.Query().Get("songid")
	songid, errors := strconv.Atoi(songidparam)
	if errors != nil {
		http.Error(Write, fmt.Sprintf("Error filtering songs: %s", errors.Error()), http.StatusInternalServerError)
		return
	}
	resp, err := Controller.songService.DeleteSongById(songid)
	if err != nil {
		http.Error(Write, fmt.Sprintf("Error filtering songs: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	Write.Header().Set("Content-Type", "application/json")
	Write.WriteHeader(http.StatusOK)
	json.NewEncoder(Write).Encode(resp)

}
