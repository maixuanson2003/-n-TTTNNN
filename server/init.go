package server

import (
	database "ten_module/Database"
	middleware "ten_module/Middleware"
	"ten_module/internal/Helper/elastichelper"
	openai "ten_module/internal/Helper/openAi"
	albumcontroller "ten_module/internal/controller/AlbumController"
	artistcontroller "ten_module/internal/controller/ArtistController"
	Authcontroller "ten_module/internal/controller/AuthController"
	collectioncontroller "ten_module/internal/controller/CollectionController"
	countrycontroller "ten_module/internal/controller/CountryController"
	historycontroller "ten_module/internal/controller/HistoryController"
	otpcontroller "ten_module/internal/controller/OtpController"
	playlistcontroller "ten_module/internal/controller/PlayListController"
	reviewcontroller "ten_module/internal/controller/ReviewController"
	songcontroller "ten_module/internal/controller/SongController"
	songtypecontroller "ten_module/internal/controller/SongTypeController"
	"ten_module/internal/controller/UserController"
	"ten_module/internal/repository"
	"ten_module/internal/service/authservice"
)

func InitSingleton() {
	middleware.InitMiddleWare()
	openai.InitGeminiClient()
	authservice.Init()
	database.Init()
	elastichelper.InitElasticHelpers()
	// Repository
	repository.InitUserRepo()

	repository.InitArtistRepository()

	repository.InitSongRepo()

	repository.InitSongTypeRepository()

	repository.InitListenHistoryRepo()

	repository.InitPlayListRepository()

	repository.InitCollectionRepostiory()

	repository.InitAlbumRepository()

	repository.InitReviewRepository()

	repository.InitCountryRepository()

	repository.InitOtpRepository()

	//Controller init
	UserController.InitService()

	Authcontroller.InitController()

	songcontroller.InitSongService()

	historycontroller.InitHistoryService()

	artistcontroller.InitArtistControll()

	playlistcontroller.InitPlayListControll()

	collectioncontroller.InitCollectionControll()

	albumcontroller.InitAlbumControll()

	reviewcontroller.InitReviewControll()

	countrycontroller.InitCountryControll()

	songtypecontroller.InitSongTypeControll()
	otpcontroller.InitOtpControll()
}
