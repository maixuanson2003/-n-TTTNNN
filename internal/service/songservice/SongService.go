package songservice

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"ten_module/internal/Config"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	helper "ten_module/internal/Helper/contentbase"
	"ten_module/internal/Helper/elastichelper"
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
		Status:       "Release",
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
func GetWeekRange() (time.Time, time.Time) {
	now := time.Now()
	weekday := int(now.Weekday())

	if weekday == 0 {
		weekday = 7
	}

	startOfWeek := now.AddDate(0, 0, -weekday+1)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	endOfWeek := startOfWeek.AddDate(0, 0, 6)
	endOfWeek = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 0, endOfWeek.Location())

	return startOfWeek, endOfWeek
}
func (songServe *SongService) GetBookTopWeek() []response.SongResponse {
	SongRepo := songServe.SongRepo
	ListSong, err := SongRepo.FindAll()
	if err != nil {
		log.Print(err)
		return nil
	}
	startWeek, endWeek := GetWeekRange()
	pairSong := make(map[*entity.Song]int32)
	for _, SongItem := range ListSong {
		count := int32(0)
		for _, listenItem := range SongItem.ListenHistory {
			if !listenItem.ListenDay.Before(startWeek) && !listenItem.ListenDay.After(endWeek) {
				count++
			}
		}
		if count > 0 {
			pairSong[&SongItem] = count
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
func (songServe *SongService) UpdateSongAlbum(AlbumId *int, songReqs []request.SongRequest, songFiles []request.SongFile) {
	album, err := songServe.AlbumRepo.GetAlbumById(*AlbumId)
	if err != nil {
		log.Print(album)
		return
	}

	err = songServe.SongRepo.DeleteSongsByAlbumId(*AlbumId)
	if err != nil {
		log.Printf("Lỗi khi xóa bài hát cũ: %v", err)
		return
	}

	album.Song = []entity.Song{}

	errs := songServe.AlbumRepo.UpdateAlbum(album, *AlbumId)
	if errs != nil {
		log.Print(errs)
		return
	}
	for index, songReq := range songReqs {
		go func() {
			songEntity := entity.Song{}

			songEntity.NameSong = songReq.NameSong
			songEntity.Description = songReq.Description
			songEntity.Point = songReq.Point
			songEntity.CountryId = songReq.CountryId
			songEntity.ReleaseDay = songReq.ReleaseDay
			songEntity.CreateDay = time.Now()
			songEntity.UpdateDay = time.Now()
			songEntity.AlbumId = AlbumId
			newSongTypes := []entity.SongType{}
			for _, typeId := range songReq.SongType {
				songType, err := songServe.SongTypeRepo.GetSongTypeById(typeId)
				if err != nil {
					log.Printf("Không tìm thấy thể loại ID %d: %v", typeId, err)
					continue
				}
				newSongTypes = append(newSongTypes, songType)
			}
			songEntity.SongType = newSongTypes

			newArtists := []entity.Artist{}
			for _, artistId := range songReq.Artist {
				artist, err := songServe.ArtistRepo.GetArtistById(artistId)
				if err != nil {
					log.Printf("Không tìm thấy nghệ sĩ ID %d: %v", artistId, err)
					continue
				}
				newArtists = append(newArtists, artist)
			}
			songEntity.Artist = newArtists
			if index < len(songFiles) && songFiles[index].File != nil {
				resourceUrl, err := Config.HandleUpLoadFile(songFiles[index].File, songReq.NameSong)
				if err != nil {
					log.Printf("Upload bài hát thất bại: %v", err)
				} else {
					songEntity.SongResource = resourceUrl
				}
			}

			errToCreate := songServe.SongRepo.CreateSong(songEntity)
			if errToCreate != nil {
				log.Print(errToCreate)
				log.Printf("update song that bai")
			} else {
				log.Printf("update song thành công")
			}

		}()
	}
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
	Song.Point = SongReq.Point
	Song.CountryId = SongReq.CountryId
	go func() {
		if SongFile.File != nil {
			errs := songServe.UpLoadSongBackgroundUpdate(Song, SongFile)
			if errs != nil {
				log.Print("some thing wrong")
			}
			log.Print("upload success")

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
func (songServe *SongService) GetAllSong(Offset int) ([]map[string]interface{}, error) {
	SongRepos := songServe.SongRepo
	ListSong, ErrorToGetListSong := SongRepos.Paginate(Offset)
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
func (songServe *SongService) GetSongById(Id int) (response.SongResponse, error) {
	SongRepos := songServe.SongRepo
	Song, ErrorToGetSong := SongRepos.GetSongById(Id)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return response.SongResponse{}, ErrorToGetSong
	}
	SongResponse := SongEntityMapToSongResponse(Song)
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

type HistoryPair struct {
	IdType int
	Amount int
}
type HistoryLike struct {
	IdType int
	Amount int
}

func TrackSongForUser(user entity.User) ([]HistoryPair, []HistoryLike, error) {
	SongUserListen := user.ListenHistory
	SongUserLike := user.Song

	TrackSongListen := make(map[int]int)
	TrackSongLike := make(map[int]int)
	ArrayHistory := []HistoryPair{}
	ArraySongLike := []HistoryLike{}
	now := time.Now()
	beginningOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	sevenDaysAgo := beginningOfToday.AddDate(0, 0, -7)
	for _, ListenHistoryItem := range SongUserListen {
		TimeUserListen := ListenHistoryItem.ListenDay
		if TimeUserListen.Before(sevenDaysAgo) || TimeUserListen.After(now) {
			continue
		}
		Song, ErrorToGetSong := SongServices.SongRepo.GetSongById(ListenHistoryItem.SongId)
		if ErrorToGetSong != nil {
			log.Print(ErrorToGetSong)
			return nil, nil, ErrorToGetSong
		}
		SongTypeUser := Song.SongType
		for _, SongTypeItem := range SongTypeUser {
			TrackSongListen[SongTypeItem.ID]++
		}
	}
	for IdSongType, value := range TrackSongListen {
		Check := HistoryPair{
			IdType: IdSongType,
			Amount: value,
		}
		ArrayHistory = append(ArrayHistory, Check)
	}
	sort.Slice(ArrayHistory, func(i, j int) bool {
		return ArrayHistory[i].Amount > ArrayHistory[j].Amount
	})
	for _, SongUserLikeItem := range SongUserLike {
		SongTypeUser, errorToGetSong := SongServices.SongRepo.GetSongById(SongUserLikeItem.ID)
		if errorToGetSong != nil {
			log.Print(errorToGetSong)
			return nil, nil, errorToGetSong
		}
		for _, SongTypeItem := range SongTypeUser.SongType {
			TrackSongLike[SongTypeItem.ID]++
		}
	}
	for IdSongType, value := range TrackSongLike {
		Check := HistoryLike{
			IdType: IdSongType,
			Amount: value,
		}
		ArraySongLike = append(ArraySongLike, Check)
	}
	fmt.Print(ArraySongLike)
	sort.Slice(ArraySongLike, func(i, j int) bool {
		return ArraySongLike[i].Amount > ArraySongLike[j].Amount
	})
	return ArrayHistory, ArraySongLike, nil

}
func GetMax(Limit int, SongLength int) int {
	if SongLength < Limit {
		return SongLength
	}
	return Limit

}
func GetTrueResult(Response []response.SongResponse) []response.SongResponse {
	checktrue := make(map[response.SongResponse]int)
	result := []response.SongResponse{}
	for _, value := range Response {
		checktrue[value]++
	}
	for key, _ := range checktrue {
		result = append(result, key)

	}
	return result

}
func (songServe *SongService) GetListSongForUser(userId string) ([]response.SongResponse, error) {
	UserRepo := songServe.UserRepo
	SongRepo := songServe.SongRepo

	SongTypeRepo := songServe.SongTypeRepo
	SongResponse := []response.SongResponse{}
	UserById, ErrorToGetUser := UserRepo.FindById(userId)
	if ErrorToGetUser != nil {
		log.Print(ErrorToGetUser)
		return nil, ErrorToGetUser
	}
	MaxListenIn7Day, MaxLike, ErrorToGet := TrackSongForUser(UserById)
	if ErrorToGet != nil {
		return nil, ErrorToGetUser
	}
	amountSongType := 0
	srcSong := rand.NewSource(time.Now().UnixNano())
	randSong := rand.New(srcSong)
	if len(MaxListenIn7Day) != 0 {
		for _, value := range MaxListenIn7Day {
			SongType, ErrorToGetType := SongTypeRepo.GetSongTypeById(value.IdType)
			if ErrorToGetType != nil {
				log.Print(ErrorToGetType)
				return nil, ErrorToGetType
			}
			SongArray := SongType.Song
			lenCheck := len(SongArray)
			fmt.Print(SongArray)
			if amountSongType == 0 {
				for i := 0; i < int(GetMax(FIRST_SONG, len(SongArray))); i++ {
					if lenCheck <= 5 {
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
					}
					if lenCheck > 5 {
						randomNumber := randSong.Intn(lenCheck)
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[randomNumber]))
					}
				}
			}
			if amountSongType == 1 {
				for i := 0; i < int(GetMax(SECOND_SONG, len(SongArray))); i++ {
					if lenCheck <= 5 {
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
					}
					if lenCheck > 5 {
						randomNumber := randSong.Intn(lenCheck)
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[randomNumber]))
					}
				}
			}
			if amountSongType == 2 {
				for i := 0; i < int(GetMax(THIRD_SONG, len(SongArray))); i++ {
					if lenCheck <= 5 {
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
					}
					if lenCheck > 5 {
						randomNumber := randSong.Intn(lenCheck)
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[randomNumber]))
					}
				}
			}
			if amountSongType > 2 {
				break
			}
			amountSongType++
		}
		return GetTrueResult(SongResponse), nil

	}
	if len(MaxLike) != 0 {
		fmt.Print("check")
		for _, value := range MaxLike {
			SongType, ErrorToGetType := SongTypeRepo.GetSongTypeById(value.IdType)
			if ErrorToGetType != nil {
				log.Print(ErrorToGetType)
				return nil, ErrorToGetType
			}
			SongArray := SongType.Song
			lenCheck := len(SongArray)
			if amountSongType == 0 {
				for i := 0; i < int(GetMax(FIRST_SONG, len(SongArray))); i++ {
					if lenCheck <= 5 {
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
					}
					if lenCheck > 5 {
						randomNumber := randSong.Intn(lenCheck)
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[randomNumber]))
					}
				}
			}
			if amountSongType == 1 {
				for i := 0; i < int(GetMax(SECOND_SONG, len(SongArray))); i++ {
					if lenCheck <= 5 {
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
					}
					if lenCheck > 5 {
						randomNumber := randSong.Intn(lenCheck)
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[randomNumber]))
					}
				}
			}
			if amountSongType == 2 {
				for i := 0; i < int(GetMax(THIRD_SONG, len(SongArray))); i++ {
					if lenCheck <= 5 {
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
					}
					if lenCheck > 5 {
						randomNumber := randSong.Intn(lenCheck)
						SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[randomNumber]))
					}
				}
			}
			if amountSongType > 2 {
				break
			}
			amountSongType++
		}
		return GetTrueResult(SongResponse), nil
	}
	fmt.Print("check2")
	Song, err := SongRepo.FindAll()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	for _, Song := range Song {
		SongResponse = append(SongResponse, SongEntityMapToSongResponse(Song))
	}
	return SongResponse, nil
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
	ListenRepos := SongServe.HisRepo
	type SongSimilarity struct {
		ID         int
		Similarity float64
	}
	vectorFeatureSong, FeatureTag, Error := helper.GetVectorFeatureForSong()
	log.Print(vectorFeatureSong)
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
		SongId = append(SongId, Item.SongId)
	}

	similarityMap := map[int]float64{}
	for _, SongIds := range SongId {
		SongItemId := SongIds
		count, err := ListenRepos.CountNumberSongId(SongIds)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		vectorSongItemId, errorToCaculate := helper.GetVectorFeatureForUser(SongItemId, FeatureTag)
		if errorToCaculate != nil {
			log.Print(errorToCaculate)
			return nil, errorToCaculate
		}
		for Id, vector := range vectorFeatureSong {
			if Id == SongItemId {
				continue
			}
			SimilarScore := helper.GetCosineSimilar(vectorSongItemId, vector)
			weightedSimilarity := SimilarScore * math.Sqrt(float64(count))
			similarityMap[Id] += weightedSimilarity
		}
	}
	similarities := []SongSimilarity{}
	for Id, Score := range similarityMap {
		similarities = append(similarities, SongSimilarity{
			ID:         Id,
			Similarity: Score,
		})
	}
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	// Lấy top 5 gợi ý (có thể điều chỉnh)
	topN := 7
	SongRecommendId := []response.SongResponse{}
	for i := 0; i < topN && i < len(similarities); i++ {
		// fmt.Print(similarities[i].Similarity, " ", similarities[i].ID)
		// fmt.Print(" ")
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
