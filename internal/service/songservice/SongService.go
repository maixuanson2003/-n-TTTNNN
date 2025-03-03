package songservice

import (
	"errors"
	"log"
	"net/http"
	"sort"
	"ten_module/internal/Config"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
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
	GetListSongForUser(userId int) ([]response.SongResponse, error)
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
		SongTypeUser := SongUserLikeItem.SongType
		for _, SongTypeItem := range SongTypeUser {
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
	sort.Slice(ArraySongLike, func(i, j int) bool {
		return ArraySongLike[i].Amount > ArraySongLike[j].Amount
	})
	return ArrayHistory, ArraySongLike, nil

}
func GetMax(limt int, Songlength int) int {
	if Songlength < limt {
		return Songlength
	}
	return limt

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
	if len(MaxListenIn7Day) != 0 {
		for _, value := range MaxListenIn7Day {
			SongType, ErrorToGetType := SongTypeRepo.GetSongTypeById(value.IdType)
			if ErrorToGetType != nil {
				log.Print(ErrorToGetType)
				return nil, ErrorToGetType
			}
			SongArray := SongType.Song
			if amountSongType == 0 {
				for i := 0; i < int(GetMax(FIRST_SONG, len(SongArray))); i++ {
					SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
				}
			}
			if amountSongType == 1 {
				for i := 0; i < int(GetMax(SECOND_SONG, len(SongArray))); i++ {
					SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
				}
			}
			if amountSongType == 2 {
				for i := 0; i < int(GetMax(THIRD_SONG, len(SongArray))); i++ {
					SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
				}
			}
			if amountSongType > 2 {
				break
			}
			amountSongType++
		}
		return SongResponse, nil

	}
	if len(MaxLike) != 0 {
		for _, value := range MaxLike {
			SongType, ErrorToGetType := SongTypeRepo.GetSongTypeById(value.IdType)
			if ErrorToGetType != nil {
				log.Print(ErrorToGetType)
				return nil, ErrorToGetType
			}
			SongArray := SongType.Song
			if amountSongType == 0 {
				for i := 0; i < int(GetMax(FIRST_SONG, len(SongArray))); i++ {
					SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
				}
			}
			if amountSongType == 1 {
				for i := 0; i < int(GetMax(SECOND_SONG, len(SongArray))); i++ {
					SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
				}
			}
			if amountSongType == 2 {
				for i := 0; i < int(GetMax(THIRD_SONG, len(SongArray))); i++ {
					SongResponse = append(SongResponse, SongEntityMapToSongResponse(SongArray[i]))
				}
			}
			if amountSongType > 2 {
				break
			}
			amountSongType++
		}
		return SongResponse, nil
	}
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
