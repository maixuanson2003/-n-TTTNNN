package songservice

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"ten_module/internal/Config"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/Helper/elastichelper"
	"ten_module/internal/repository"
	"time"
)

type SongService struct {
	UserRepo     *repository.UserRepository
	SongRepo     *repository.SongRepository
	SongTypeRepo *repository.SongTypeRepository
	ArtistRepo   *repository.ArtistRepository
}
type SongServiceInterface interface {
	GetSongById(Id int) (response.SongResponse, error)
	GetAllSong() ([]response.SongResponse, error)
	CreateNewSong(SongReq request.SongRequest, SongFile request.SongFile) (MessageResponse, error)
	DownLoadSong(Id int) (SongDownload, error)
	GetListSongForUser(userId string) ([]response.SongResponse, error)
	UpdateSong(SongReq request.SongRequest, Id int) (MessageResponse, error)
	UserLikeSong(SongId int, UserId string) (MessageResponse, error)
	SearchSong(Keyword string) ([]response.SongResponse, error)
	FilterSong()
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

func InitSongService() {
	SongServices = &SongService{
		UserRepo:     repository.UserRepo,
		SongRepo:     repository.SongRepo,
		SongTypeRepo: repository.SongTypeRepo,
		ArtistRepo:   repository.ArtistRepo,
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
	resourceSong, err := Config.HandleUpLoadFile(SongFile.File, SongReq.NameSong)
	if SongReq.NameSong == "" {
		return MessageResponse{
			Message: "Failed to create",
			Status:  "Failed",
		}, errors.New("name song is empty")
	}
	if err != nil {
		return MessageResponse{
			Message: "Failed to create",
			Status:  "Failed",
		}, err
	}
	SongEntity := SongReqMapToSongEntity(SongReq, resourceSong, ListSongType, ListArtist)
	errorToCreateSong := songServe.SongRepo.CreateSong(SongEntity)
	if errorToCreateSong != nil {
		return MessageResponse{
			Message: "failed to create song",
			Status:  "failed",
		}, errorToCreateSong
	}
	return MessageResponse{
		Message: "Success to create song",
		Status:  "Success",
	}, nil

}
func (songServe *SongService) GetAllSong() ([]response.SongResponse, error) {
	SongRepos := songServe.SongRepo
	ListSong, ErrorToGetListSong := SongRepos.FindAll()
	if ErrorToGetListSong != nil {
		log.Print(ErrorToGetListSong)
		return nil, ErrorToGetListSong
	}
	ListSongResponse := []response.SongResponse{}
	for _, SongItem := range ListSong {
		ListSongResponse = append(ListSongResponse, SongEntityMapToSongResponse(SongItem))
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
