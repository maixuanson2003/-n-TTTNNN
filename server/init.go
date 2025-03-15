package server

import (
	database "ten_module/Database"
	middleware "ten_module/Middleware"
	albumcontroller "ten_module/internal/controller/AlbumController"
	artistcontroller "ten_module/internal/controller/ArtistController"
	Authcontroller "ten_module/internal/controller/AuthController"
	collectioncontroller "ten_module/internal/controller/CollectionController"
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
	// Repository
	repository.InitUserRepo()
	repository.InitArtistRepository()
	repository.InitSongRepo()
	repository.InitSongTypeRepository()
	repository.InitListenHistoryRepo()
	repository.InitPlayListRepository()
	repository.InitCollectionRepostiory()
	repository.InitAlbumRepository()
	UserController.InitService()
	//Controller init
	Authcontroller.InitController()
	songcontroller.InitSongService()
	historycontroller.InitHistoryService()
	artistcontroller.InitArtistControll()
	playlistcontroller.InitPlayListControll()
	collectioncontroller.InitCollectionControll()
	albumcontroller.InitAlbumControll()
}
