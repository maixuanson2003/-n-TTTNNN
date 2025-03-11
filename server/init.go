package server

import (
	database "ten_module/Database"
	middleware "ten_module/Middleware"
	artistcontroller "ten_module/internal/controller/ArtistController"
	Authcontroller "ten_module/internal/controller/AuthController"
	historycontroller "ten_module/internal/controller/HistoryController"
	playlistcontroller "ten_module/internal/controller/PlayListController"
	songcontroller "ten_module/internal/controller/SongController"
	"ten_module/internal/controller/UserController"
	"ten_module/internal/repository"
	"ten_module/internal/service/authservice"
)

func InitSingleton() {
	middleware.InitMiddleWare()
	authservice.Init()
	database.Init()
	repository.InitUserRepo()
	repository.InitArtistRepository()
	repository.InitSongRepo()
	repository.InitSongTypeRepository()
	repository.InitListenHistoryRepo()
	repository.InitPlayListRepository()
	repository.InitCollectionRepostiory()
	UserController.InitService()
	Authcontroller.InitController()
	songcontroller.InitSongService()
	historycontroller.InitHistoryService()
	artistcontroller.InitArtistControll()
	playlistcontroller.InitPlayListControll()
}
