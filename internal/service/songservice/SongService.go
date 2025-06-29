package songservice

import (
	"errors"
	"fmt"
	"log"

	"net/http"
	"sort"
	"ten_module/internal/Config"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	helper "ten_module/internal/Helper/contentbase"
	"ten_module/internal/Helper/elastichelper"
	gemini "ten_module/internal/Helper/openAi"
	"ten_module/internal/repository"
	"time"
)

type SongService struct {
	UserRepo     *repository.UserRepository
	SongRepo     *repository.SongRepository
	SongTypeRepo *repository.SongTypeRepository
	ArtistRepo   *repository.ArtistRepository
	HisRepo      *repository.ListenHistoryRepo
	AlbumRepo    *repository.AlbumRepository
}
type SongServiceInterface interface {
	GetSongById(Id int) (response.SongResponse, error)
	GetAllSong(Offset int) ([]map[string]interface{}, error)
	CreateNewSong(SongReq request.SongRequest, SongFile request.SongFile) (MessageResponse, error)
	DownLoadSong(Id int) (SongDownload, error)
	GetListSongForUser(userId string) ([]response.SongResponse, error)
	UpdateSong(SongReq request.SongRequest, Id int) (MessageResponse, error)
	UserLikeSong(SongId int, UserId string) (MessageResponse, error)
	SearchSong(Keyword string) ([]response.SongResponse, error)
	GetSongForUser(userId string) ([]response.SongResponse, error)
	FilterSong(ArtistId []int, TypeId []int) ([]map[string]interface{}, error)
	SearchSongByKeyWord(keyWord string) ([]map[string]interface{}, error)
}
type MessageResponse struct {
	Message string
	Status  string
}

const (
	FIRST_SONG  = 4
	SECOND_SONG = 3
	THIRD_SONG  = 2
)

var SongServices *SongService
var VectorFeature []string
var VectorAllSong map[int][]int16

func InitSongService() {
	SongServices = &SongService{
		UserRepo:     repository.UserRepo,
		SongRepo:     repository.SongRepo,
		SongTypeRepo: repository.SongTypeRepo,
		ArtistRepo:   repository.ArtistRepo,
		HisRepo:      repository.ListenRepo,
		AlbumRepo:    repository.AlbumRepo,
	}
}
func SongReqMapToSongEntity(SongReq request.SongRequest, resource string, ListSongType []entity.SongType, ListArtist []entity.Artist) entity.Song {
	return entity.Song{
		NameSong:     SongReq.NameSong,
		Description:  SongReq.Description,
		ReleaseDay:   time.Now(),
		CreateDay:    time.Now(),
		UpdateDay:    time.Now(),
		Point:        SongReq.Point,
		LikeAmount:   0,
		Status:       SongReq.Status,
		CountryId:    SongReq.CountryId,
		ListenAmout:  0,
		SongResource: resource,
		SongType:     ListSongType,
		Artist:       ListArtist,
	}
}
func SongEntityMapToSongResponse(Song entity.Song) response.SongResponse {
	return response.SongResponse{
		ID:           Song.ID,
		NameSong:     Song.NameSong,
		Description:  Song.Description,
		ReleaseDay:   Song.ReleaseDay,
		CreateDay:    Song.CreateDay,
		UpdateDay:    Song.UpdateDay,
		Point:        Song.Point,
		LikeAmount:   Song.LikeAmount,
		Status:       Song.Status,
		CountryId:    Song.CountryId,
		ListenAmout:  Song.ListenAmout,
		AlbumId:      Song.AlbumId,
		SongResource: Song.SongResource,
	}
}
func SongEntityMapToSongResponseAlbum(Song entity.Song, countryrep *repository.CountryRepository) response.SongResponseAlbum {

	songArtistResponses := []response.ArtistResponse{}

	for _, artist := range Song.Artist {
		country, err := countryrep.GetCountryById(artist.CountryId)
		if err != nil {
			log.Printf("Lỗi khi lấy quốc gia của nghệ sĩ ID %d: %v", artist.ID, err)
			continue
		}
		songArtistResponses = append(songArtistResponses, response.ArtistResponse{
			ID:          artist.ID,
			Name:        artist.Name,
			BirthDay:    artist.BirthDay,
			Description: artist.Description,
			Country:     country.CountryName,
		})
	}

	// Map thể loại của bài hát
	songTypeResponses := []response.SongTypeResponse{}
	for _, songType := range Song.SongType {
		songTypeResponses = append(songTypeResponses, response.SongTypeResponse{Id: songType.ID, Type: songType.Type})
	}
	return response.SongResponseAlbum{
		ID:           Song.ID,
		NameSong:     Song.NameSong,
		Description:  Song.Description,
		ReleaseDay:   Song.ReleaseDay,
		CreateDay:    Song.CreateDay,
		UpdateDay:    Song.UpdateDay,
		Point:        Song.Point,
		LikeAmount:   Song.LikeAmount,
		Status:       Song.Status,
		CountryId:    Song.CountryId,
		ListenAmout:  Song.ListenAmout,
		AlbumId:      Song.AlbumId,
		SongResource: Song.SongResource,
		Artist:       songArtistResponses,
		SongType:     songTypeResponses,
	}
}
func GetTimeRange(rangeType string) (time.Time, time.Time) {
	now := time.Now()
	var start, end time.Time

	switch rangeType {
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = now.AddDate(0, 0, -weekday+1)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		end = start.AddDate(0, 0, 6)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, -1)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())
	default:
		// mặc định là tuần
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = now.AddDate(0, 0, -weekday+1)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		end = start.AddDate(0, 0, 6)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())
	}

	return start, end
}

type SongChartData struct {
	Name         string `json:"name"`
	ListenPerDay [7]int `json:"listen_per_day"` // thứ 2 -> chủ nhật
}

func (songServe *SongService) GetDataFromPromb(query *gemini.MusicQuery) []response.SongResponse {
	log.Print(query)
	song, errs := songServe.SongRepo.RecommendSongs(query.Genre, query.Artist, query.Country, query.Keywords, query.TimeRange, query.SortBy)
	if errs != nil {
		return nil
	}
	SongResponse := []response.SongResponse{}

	for i := 0; i < len(song); i++ {
		SongResponse = append(SongResponse, SongEntityMapToSongResponse(song[i]))
	}
	return SongResponse

}
func (songServe *SongService) GetWeeklyChartDataPerDay(topN int) []SongChartData {
	startWeek, endWeek := GetTimeRange("week")
	SongRepo := songServe.SongRepo
	ListSong, err := SongRepo.FindAll()
	if err != nil {
		log.Print(err)
		return nil
	}

	type songWithListen struct {
		Song         entity.Song
		Total        int
		ListenPerDay [7]int
	}

	var songsData []songWithListen

	for i := range ListSong {
		var listenPerDay [7]int
		total := 0
		for _, hist := range ListSong[i].ListenHistory {
			if !hist.ListenDay.Before(startWeek) && !hist.ListenDay.After(endWeek) {
				dayIdx := int(hist.ListenDay.Weekday())
				if dayIdx == 0 {
					dayIdx = 6
				} else {
					dayIdx = dayIdx - 1
				}
				listenPerDay[dayIdx]++
				total++
			}
		}
		// Chỉ lấy các bài có lượt nghe > 0
		if total > 0 {
			songsData = append(songsData, songWithListen{
				Song:         ListSong[i],
				Total:        total,
				ListenPerDay: listenPerDay,
			})
		}
	}

	sort.Slice(songsData, func(i, j int) bool {
		return songsData[i].Total > songsData[j].Total
	})

	if len(songsData) > topN {
		songsData = songsData[:topN]
	}

	var chart []SongChartData
	for _, item := range songsData {
		chart = append(chart, SongChartData{
			Name:         item.Song.NameSong,
			ListenPerDay: item.ListenPerDay,
		})
	}

	return chart
}
func (songServe *SongService) GetBookTopRange(rangeType string) []response.SongResponse {
	SongRepo := songServe.SongRepo
	ListSong, err := SongRepo.FindAll()
	if err != nil {
		log.Print(err)
		return nil
	}
	startRange, endRange := GetTimeRange(rangeType)
	pairSong := make(map[*entity.Song]int32)
	for i := range ListSong {
		SongItem := &ListSong[i]
		count := int32(0)
		for _, listenItem := range SongItem.ListenHistory {
			if !listenItem.ListenDay.Before(startRange) && !listenItem.ListenDay.After(endRange) {
				count++
			}
		}
		if count > 0 {
			pairSong[SongItem] = count
		}
	}
	type songSlice struct {
		Song   entity.Song
		amount int32
	}
	arraySong := []songSlice{}
	for Song, count := range pairSong {
		arraySong = append(arraySong, songSlice{
			Song:   *Song,
			amount: count,
		})
	}
	sort.Slice(arraySong, func(i, j int) bool {
		return arraySong[i].amount > arraySong[j].amount
	})
	SongResponse := []response.SongResponse{}
	limit := 5
	if len(arraySong) < 5 {
		limit = len(arraySong)
	}
	for i := 0; i < limit; i++ {
		SongResponse = append(SongResponse, SongEntityMapToSongResponse(arraySong[i].Song))
	}

	return SongResponse
}
func (songServe *SongService) CreateNewSong(SongReq request.SongRequest, SongFile request.SongFile) (MessageResponse, error) {
	ListSongType := []entity.SongType{}
	ListArtist := []entity.Artist{}
	for _, IdSongType := range SongReq.SongType {
		SongType, err := songServe.SongTypeRepo.GetSongTypeById(IdSongType)
		if err != nil {
			log.Print(err)
			return MessageResponse{}, err
		}
		ListSongType = append(ListSongType, SongType)
	}
	for _, IdArtist := range SongReq.Artist {
		Artist, err := songServe.ArtistRepo.GetArtistById(IdArtist)
		if err != nil {
			log.Print(err)
			return MessageResponse{}, err
		}
		ListArtist = append(ListArtist, Artist)
	}
	go func() {
		errs := songServe.UpLoadSongBackground(SongReq, SongFile, ListSongType, ListArtist)
		if errs != nil {
			log.Print("some thing wrong")
		}
		log.Print("upload success")

	}()
	return MessageResponse{
		Message: "Success to create song",
		Status:  "Success",
	}, nil

}
func (songServe *SongService) UpLoadSongBackground(SongReq request.SongRequest, SongFile request.SongFile, ListSongType []entity.SongType, ListArtist []entity.Artist) error {
	errorRes := errors.New("check")
	resourceSong, err := Config.HandleUpLoadFile(SongFile.File, SongReq.NameSong)
	if SongReq.NameSong == "" {
		errorRes = errors.New("require name song")
	}
	if err != nil {
		errorRes = err
	}
	SongEntity := SongReqMapToSongEntity(SongReq, resourceSong, ListSongType, ListArtist)
	errorToCreateSong := songServe.SongRepo.CreateSong(SongEntity)
	if errorToCreateSong != nil {
		errorRes = errorToCreateSong
	}
	return errorRes
}
func (songServe *SongService) UpLoadSongBackgroundUpdate(songEntity entity.Song, SongFile request.SongFile) error {
	errorRes := errors.New("check")
	resourceSong, err := Config.HandleUpLoadFile(SongFile.File, songEntity.NameSong)
	if songEntity.NameSong == "" {
		errorRes = errors.New("require name song")
	}
	if err != nil {
		errorRes = err
	}
	songEntity.SongResource = resourceSong
	errorToUpdateSong := songServe.SongRepo.UpdateSong(songEntity, songEntity.ID)
	if errorToUpdateSong != nil {
		errorRes = errorToUpdateSong
	}
	return errorRes
}
func (songServe *SongService) UpdateSong(SongReq request.SongRequest, Id int, SongFile request.SongFile) (MessageResponse, error) {
	Song, errToGetSong := songServe.SongRepo.GetSongById(Id)
	if errToGetSong != nil {
		return MessageResponse{
			Message: "failed",
			Status:  "Failed",
		}, errToGetSong
	}
	ListSongType := []entity.SongType{}
	ListArtist := []entity.Artist{}
	for _, IdSongType := range SongReq.SongType {
		SongType, err := songServe.SongTypeRepo.GetSongTypeById(IdSongType)
		if err != nil {
			log.Print(err)
			return MessageResponse{}, err
		}
		ListSongType = append(ListSongType, SongType)
	}
	for _, IdArtist := range SongReq.Artist {
		Artist, err := songServe.ArtistRepo.GetArtistById(IdArtist)
		if err != nil {
			log.Print(err)
			return MessageResponse{}, err
		}
		ListArtist = append(ListArtist, Artist)
	}
	Song.SongType = ListSongType
	Song.Artist = ListArtist
	Song.NameSong = SongReq.NameSong
	Song.Description = SongReq.Description
	Song.Status = SongReq.Status
	Song.Point = SongReq.Point
	Song.CountryId = SongReq.CountryId
	go func() {
		if SongFile.File != nil {
			if err := songServe.UpLoadSongBackgroundUpdate(Song, SongFile); err != nil {
				log.Print("Có lỗi khi upload:", err)
			} else {
				log.Print("Upload thành công")
			}
		} else {
			log.Print("Không có file (nil hoặc typed-nil)")
		}

	}()
	err := songServe.SongRepo.UpdateSong(Song, Id)
	if err != nil {
		return MessageResponse{
			Message: "failed",
			Status:  "Failed",
		}, err
	}
	return MessageResponse{
		Message: "success",
		Status:  "Success",
	}, err

}
func (songServe *SongService) GetListSong() ([]map[string]interface{}, error) {
	SongRepos := songServe.SongRepo
	ListSong, ErrorToGetListSong := SongRepos.FindAll()
	if ErrorToGetListSong != nil {
		log.Print(ErrorToGetListSong)
		return nil, ErrorToGetListSong
	}
	ListSongResponse := []map[string]interface{}{}
	for _, SongItem := range ListSong {
		SongResponseItem := SongEntityMapToSongResponse(SongItem)
		Aritst := SongItem.Artist
		ArtistForSong := []map[string]interface{}{}
		for _, item := range Aritst {
			ArtistRes := map[string]interface{}{
				"id":          item.ID,
				"name":        item.Name,
				"description": item.Description,
			}
			ArtistForSong = append(ArtistForSong, ArtistRes)
		}
		Songs := map[string]interface{}{
			"SongData": SongResponseItem,
			"artist":   ArtistForSong,
		}
		ListSongResponse = append(ListSongResponse, Songs)

	}
	return ListSongResponse, nil
}
func (songServe *SongService) GetAllSong() ([]map[string]interface{}, error) {
	SongRepos := songServe.SongRepo
	ListSong, ErrorToGetListSong := SongRepos.FindAll()
	if ErrorToGetListSong != nil {
		log.Print(ErrorToGetListSong)
		return nil, ErrorToGetListSong
	}
	ListSongResponse := []map[string]interface{}{}
	for _, SongItem := range ListSong {
		log.Print(SongItem.Status)
		if SongItem.Status == "public" {
			SongResponseItem := SongEntityMapToSongResponse(SongItem)
			Aritst := SongItem.Artist
			ArtistForSong := []map[string]interface{}{}
			for _, item := range Aritst {
				ArtistRes := map[string]interface{}{
					"id":          item.ID,
					"name":        item.Name,
					"description": item.Description,
				}
				ArtistForSong = append(ArtistForSong, ArtistRes)
			}
			Songs := map[string]interface{}{
				"SongData": SongResponseItem,
				"artist":   ArtistForSong,
			}
			ListSongResponse = append(ListSongResponse, Songs)

		}

	}
	return ListSongResponse, nil
}
func (songServe *SongService) GetSongByUserId(userid string) ([]response.SongResponseAlbum, error) {
	UserRepo := songServe.UserRepo
	User, ErrorToGetUser := UserRepo.FindById(userid)
	if ErrorToGetUser != nil {
		log.Print(ErrorToGetUser)
		return nil, ErrorToGetUser
	}
	Song := User.Song
	response := []response.SongResponseAlbum{}
	for _, item := range Song {
		response = append(response, SongEntityMapToSongResponseAlbum(item, repository.CountryRepo))
	}
	return response, nil
}
func (songServe *SongService) GetSongById(Id int) (response.SongResponseAlbum, error) {
	SongRepos := songServe.SongRepo
	Song, ErrorToGetSong := SongRepos.GetSongById(Id)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return response.SongResponseAlbum{}, ErrorToGetSong
	}
	SongResponse := SongEntityMapToSongResponseAlbum(Song, repository.CountryRepo)
	return SongResponse, nil
}

type SongDownload struct {
	Resp     *http.Response
	NameSong string
}

func (songServe *SongService) DownLoadSong(Id int) (SongDownload, error) {
	SongRepos := songServe.SongRepo
	Song, ErrorToGetSong := SongRepos.GetSongById(Id)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return SongDownload{}, ErrorToGetSong
	}
	resp, errorToGetSongAudio := Config.HandleDownLoadFile(Song.NameSong, "video")
	if errorToGetSongAudio != nil {
		log.Print(errorToGetSongAudio)
		return SongDownload{}, errorToGetSongAudio
	}
	return SongDownload{
		Resp:     resp,
		NameSong: Song.NameSong,
	}, nil
}

func (songServe *SongService) UserLikeSong(SongId int, UserId string) (MessageResponse, error) {
	UserRepo := songServe.UserRepo
	SongRepo := songServe.SongRepo
	User, ErrorToGetUser := UserRepo.FindById(UserId)
	Song, ErrorToGetSong := SongRepo.GetSongById(SongId)
	if ErrorToGetUser != nil {
		return MessageResponse{
			Message: "fail",
			Status:  "failed",
		}, ErrorToGetUser
	}
	if ErrorToGetSong != nil {
		return MessageResponse{
			Message: "fail",
			Status:  "failed",
		}, ErrorToGetSong
	}
	User.Song = append(User.Song, Song)
	Song.LikeAmount += 1
	ErrorUpdateSong := SongRepo.UpdateSong(Song, Song.ID)
	if ErrorUpdateSong != nil {
		log.Print(ErrorUpdateSong)
		return MessageResponse{
			Message: "faile",
			Status:  "failed",
		}, ErrorUpdateSong
	}
	ErrorToUpdate := UserRepo.Update(User, UserId)
	if ErrorToUpdate != nil {
		return MessageResponse{
			Message: "fail",
			Status:  "failed",
		}, ErrorToUpdate
	}
	return MessageResponse{
		Message: "Update Success",
		Status:  "Success",
	}, nil
}
func (songServe *SongService) UserDislikeSong(songID int, userID string) (MessageResponse, error) {
	userRepo := songServe.UserRepo
	songRepo := songServe.SongRepo
	user, errUser := userRepo.FindById(userID)
	song, errSong := songRepo.GetSongById(songID)

	if errUser != nil || errSong != nil {
		return MessageResponse{Message: "fail", Status: "failed"},
			fmt.Errorf("get user err: %v | get song err: %v", errUser, errSong)
	}
	found := false
	newList := make([]entity.Song, 0, len(user.Song))
	for _, s := range user.Song {
		if s.ID == songID {
			found = true
			continue
		}
		newList = append(newList, s)
	}
	if !found {
		return MessageResponse{
			Message: "User chưa like bài hát này",
			Status:  "failed",
		}, nil
	}
	user.Song = newList

	if song.LikeAmount > 0 {
		song.LikeAmount--
	}

	if err := songRepo.UpdateSong(song, song.ID); err != nil {
		return MessageResponse{Message: "fail", Status: "failed"}, err
	}
	if err := userRepo.DeleteSongLike(userID, songID); err != nil {
		return MessageResponse{Message: "fail", Status: "failed"}, err
	}

	return MessageResponse{
		Message: "Update Success",
		Status:  "Success",
	}, nil
}
func (SongServe *SongService) SearchSong(Keyword string) ([]response.SongResponse, error) {
	Elastic := elastichelper.ElasticHelpers
	SongResponse, errorToSearchSong := Elastic.SearchSong(Keyword)
	if errorToSearchSong != nil {
		log.Print(errorToSearchSong)
		return nil, errorToSearchSong
	}
	return SongResponse, nil
}
func (SongServe *SongService) GetSongForUser(userId string) ([]response.SongResponse, error) {
	UserRepo := SongServe.UserRepo
	type SongSimilarity struct {
		ID         int
		Similarity float64
	}
	vectorFeatureSong, FeatureTag, Error := helper.GetVectorFeatureForSong()
	if Error != nil {
		log.Print(Error)
		return nil, Error
	}
	UserItem, ErrorToGetUser := UserRepo.FindById(userId)
	if ErrorToGetUser != nil {
		log.Print(ErrorToGetUser)
		return nil, ErrorToGetUser
	}
	SongId := []int{}
	ListenHistory := UserItem.ListenHistory

	for _, Item := range ListenHistory {
		SongId = append(SongId, Item.ID)
	}
	userVector, errs := helper.GetUserProfile(userId, FeatureTag)
	log.Print("debug")
	log.Print(FeatureTag)
	log.Print(userVector)
	log.Print("debug")
	if errs != nil {
		return nil, errs
	}

	similarityMap := map[int]float64{}
	for _, SongIds := range SongId {
		SongItemId := SongIds
		for Id, vector := range vectorFeatureSong {
			if Id == SongItemId {
				continue
			}
			SimilarScore := helper.GetCosineSimilar(userVector, vector)
			similarityMap[Id] = SimilarScore
		}
	}
	similarities := []SongSimilarity{}
	for Id, Score := range similarityMap {
		similarities = append(similarities, SongSimilarity{
			ID:         Id,
			Similarity: Score,
		})
	}
	log.Print(similarities)
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})
	topN := 7
	SongRecommendId := []response.SongResponse{}
	for i := 0; i < topN && i < len(similarities); i++ {
		SongEntity, errorToGetSong := SongServe.SongRepo.GetSongById(similarities[i].ID)
		if errorToGetSong != nil {
			log.Print(errorToGetSong)
			return nil, errorToGetSong
		}
		SongResponse := SongEntityMapToSongResponse(SongEntity)
		SongRecommendId = append(SongRecommendId, SongResponse)
	}
	return SongRecommendId, nil
}
func (SongServe *SongService) GetSongForUserV2(userID string, topK int, topN int) ([]response.SongResponse, error) {
	user, err := SongServe.UserRepo.FindById(userID)
	if err != nil {
		return nil, err
	}
	type SongSimilarity struct {
		ID         int
		Score      float64
		Count      int
		Similarity float64
	}

	songVectors, _, err := helper.GetVectorFeatureForSong()
	if err != nil {
		return nil, err
	}

	listened := map[int]bool{}
	for _, item := range user.ListenHistory {
		listened[item.SongId] = true
	}
	scoreMap := map[int]*SongSimilarity{}
	cutoff := time.Now().AddDate(0, 0, -5)
	for _, item := range user.ListenHistory {
		if item.ListenDay.Before(cutoff) {
			continue
		}
		sourceID := item.SongId
		sourceVec := songVectors[sourceID]

		candidates := []SongSimilarity{}
		for songID, targetVec := range songVectors {
			// if songID == sourceID || listened[songID] {
			// 	continue
			// }
			sim := helper.GetCosineSimilar(sourceVec, targetVec)
			if sim > 0 {
				candidates = append(candidates, SongSimilarity{ID: songID, Similarity: sim})
			}
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Similarity > candidates[j].Similarity
		})

		for i := 0; i < topK && i < len(candidates); i++ {
			songID := candidates[i].ID
			if _, ok := scoreMap[songID]; !ok {
				scoreMap[songID] = &SongSimilarity{ID: songID}
			}
			scoreMap[songID].Count++
			scoreMap[songID].Similarity += candidates[i].Similarity
		}
	}

	final := []SongSimilarity{}
	for _, sim := range scoreMap {
		final = append(final, *sim)
	}

	sort.Slice(final, func(i, j int) bool {
		if final[i].Count == final[j].Count {
			return final[i].Similarity > final[j].Similarity
		}
		return final[i].Count > final[j].Count
	})

	results := []response.SongResponse{}
	for i := 0; i < topN && i < len(final); i++ {
		song, err := SongServe.SongRepo.GetSongById(final[i].ID)
		if err != nil {
			log.Print(err)
			continue
		}
		results = append(results, SongEntityMapToSongResponse(song))
	}
	if len(results) == 0 {
		log.Print("sss")
		return SongServe.GetBookTopRange("month"), nil
	}
	log.Print(len(results))

	return results, nil
}
func (SongServe *SongService) GetSimilarSongs(songId int) ([]response.SongResponse, error) {
	type SongSimilarity struct {
		ID         int
		Similarity float64
	}

	vectorFeatureSong, featureTags, err := helper.GetVectorFeatureForSong()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	vectorTargetSong, err := helper.GetVectorFeatureForUser(songId, featureTags)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	similarityMap := map[int]float64{}
	for id, vector := range vectorFeatureSong {
		if id == songId {
			continue // bỏ qua chính nó
		}
		score := helper.GetCosineSimilar(vectorTargetSong, vector)
		similarityMap[id] = score
	}
	similarities := []SongSimilarity{}
	for id, score := range similarityMap {
		similarities = append(similarities, SongSimilarity{
			ID:         id,
			Similarity: score,
		})
	}
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})
	topN := 7
	result := []response.SongResponse{}
	for i := 0; i < topN && i < len(similarities); i++ {
		songEntity, err := SongServe.SongRepo.GetSongById(similarities[i].ID)
		if err != nil {
			log.Print(err)
			continue // hoặc return nil, err nếu bạn muốn fail sớm
		}
		result = append(result, SongEntityMapToSongResponse(songEntity))
	}

	return result, nil
}

func (SongServe *SongService) SearchSongByKeyWord(keyWord string) ([]map[string]interface{}, error) {
	SongRepo := SongServe.SongRepo
	SongResult, errorToSearchSong := SongRepo.SearchSongByKey(keyWord)
	if errorToSearchSong != nil {
		log.Print(errorToSearchSong)
		return nil, errorToSearchSong
	}
	ListSongResponse := []map[string]interface{}{}
	for _, SongItem := range SongResult {
		SongResponseItem := SongEntityMapToSongResponse(SongItem)
		Aritst := SongItem.Artist
		ArtistForSong := []map[string]interface{}{}
		for _, item := range Aritst {
			ArtistRes := map[string]interface{}{
				"id":          item.ID,
				"name":        item.Name,
				"description": item.Description,
			}
			ArtistForSong = append(ArtistForSong, ArtistRes)
		}
		Songs := map[string]interface{}{
			"SongData": SongResponseItem,
			"artist":   ArtistForSong,
		}
		ListSongResponse = append(ListSongResponse, Songs)

	}
	return ListSongResponse, nil

}
func (SongServe *SongService) FilterSong(ArtistId []int, TypeId []int) ([]map[string]interface{}, error) {
	SongRepo := SongServe.SongRepo
	SongResult, errorToSearchSong := SongRepo.FilterSong(ArtistId, TypeId)
	if errorToSearchSong != nil {
		log.Print(errorToSearchSong)
		return nil, errorToSearchSong
	}
	ListSongResponse := []map[string]interface{}{}
	for _, SongItem := range SongResult {
		SongResponseItem := SongEntityMapToSongResponse(SongItem)
		Aritst := SongItem.Artist
		ArtistForSong := []map[string]interface{}{}
		for _, item := range Aritst {
			ArtistRes := map[string]interface{}{
				"id":          item.ID,
				"name":        item.Name,
				"description": item.Description,
			}
			ArtistForSong = append(ArtistForSong, ArtistRes)
		}
		Songs := map[string]interface{}{
			"SongData": SongResponseItem,
			"artist":   ArtistForSong,
		}
		ListSongResponse = append(ListSongResponse, Songs)

	}
	return ListSongResponse, nil

}
func (SongServe *SongService) DeleteSongById(id int) (MessageResponse, error) {
	SongRepo := SongServe.SongRepo
	err := SongRepo.DeleteSongById(id)
	if err != nil {
		return MessageResponse{
			Message: "failed",
			Status:  "Success",
		}, err
	}
	return MessageResponse{
		Message: "success",
		Status:  "Success",
	}, nil
}
