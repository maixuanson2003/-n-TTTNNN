// package server

// import (
// 	"net/http"
// 	database "ten_module/Database"
// 	albumcontroller "ten_module/internal/controller/AlbumController"
// 	artistcontroller "ten_module/internal/controller/ArtistController"
// 	Authcontroller "ten_module/internal/controller/AuthController"
// 	collectioncontroller "ten_module/internal/controller/CollectionController"
// 	countrycontroller "ten_module/internal/controller/CountryController"
// 	historycontroller "ten_module/internal/controller/HistoryController"
// 	otpcontroller "ten_module/internal/controller/OtpController"
// 	playlistcontroller "ten_module/internal/controller/PlayListController"
// 	reviewcontroller "ten_module/internal/controller/ReviewController"
// 	songcontroller "ten_module/internal/controller/SongController"
// 	songtypecontroller "ten_module/internal/controller/SongTypeController"
// 	"ten_module/internal/controller/UserController"

// 	"github.com/gorilla/mux"
// 	"github.com/rs/cors"
// 	"gorm.io/gorm"
// )

// type Server struct {
// 	Address  string
// 	Database *gorm.DB
// }

// func GetNewServer(address string, database *gorm.DB) Server {
// 	return Server{
// 		Address:  address,
// 		Database: database,
// 	}
// }
// func (server *Server) Run(address *string, databases *gorm.DB) {
// 	router := mux.NewRouter()
// 	database.MigrateDB(database.Database)
// 	go func() {
// 		database.RunSQLFile(database.Database)
// 	}()

// 	mainRouter := router.PathPrefix("/api").Subrouter()
// 	Cors := cors.New(cors.Options{
// 		AllowedOrigins:   []string{"http://localhost:3000"},
// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowCredentials: true,
// 		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "Accept"},
// 		Debug:            true,
// 	})

//		handler := Cors.Handler(router)
//		mainRouter.Use(Cors.Handler)
//		// register route
//		userController := UserController.UserControll
//		userController.RegisterRoute(mainRouter)
//		// register route
//		authController := Authcontroller.AuthControllers
//		authController.RegisterRoute(mainRouter)
//		//song route
//		songController := songcontroller.SongControllers
//		songController.RegisterRoute(mainRouter)
//		// history route
//		HistoryController := historycontroller.HistoryControllers
//		HistoryController.RegisterRoute(mainRouter)
//		//artist route
//		ArtistController := artistcontroller.ArtistControll
//		ArtistController.RegisterRoute(mainRouter)
//		//playlist route
//		PlayListController := playlistcontroller.PlayListControll
//		PlayListController.RegisterRoute(mainRouter)
//		//collection route
//		CollectionController := collectioncontroller.CollectionControll
//		CollectionController.RegisterRoute(mainRouter)
//		//album route
//		AlbumController := albumcontroller.AlbumControll
//		AlbumController.RegisterRoute(mainRouter)
//		//review route
//		ReviewController := reviewcontroller.ReviewControll
//		ReviewController.RegisterRoute(mainRouter)
//		// songtype route
//		SongTypeController := songtypecontroller.SongTypeControll
//		SongTypeController.RegisterRoute(mainRouter)
//		// country route
//		CountryController := countrycontroller.CountryControll
//		CountryController.RegisterRoute(mainRouter)
//		// otp route
//		OtpController := otpcontroller.OtpControll
//		OtpController.RegisterRoute(mainRouter)
//		fs := http.FileServer(http.Dir("C:/Users/DPC/Desktop/MusicMp4/internal/music"))
//		router.PathPrefix("/music/").Handler(http.StripPrefix("/music/", fs))
//		http.ListenAndServe(*address, handler)
//	}
package server

import (
	"net/http"
	database "ten_module/Database"
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

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

type Server struct {
	Address  string
	Database *gorm.DB
}

func GetNewServer(address string, database *gorm.DB) Server {
	return Server{
		Address:  address,
		Database: database,
	}
}

func (server *Server) Run(address *string, databases *gorm.DB) {
	router := mux.NewRouter()
	database.MigrateDB(database.Database)

	// Run SQL in background
	go func() {
		database.RunSQLFile(database.Database)
		database.Seed(database.Database)
	}()

	// Main API router
	mainRouter := router.PathPrefix("/api").Subrouter()

	// Cấu hình CORS
	Cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "Accept"},
		AllowCredentials: true,
		Debug:            true,
	})

	// ✅ Đúng cách để áp dụng CORS
	handler := Cors.Handler(router)

	// ❌ Không dùng dòng này nữa:
	// mainRouter.Use(Cors.Handler)

	// Register route cho từng controller
	UserController := UserController.UserControll
	UserController.RegisterRoute(mainRouter)

	AuthController := Authcontroller.AuthControllers
	AuthController.RegisterRoute(mainRouter)

	SongController := songcontroller.SongControllers
	SongController.RegisterRoute(mainRouter)

	HistoryController := historycontroller.HistoryControllers
	HistoryController.RegisterRoute(mainRouter)

	ArtistController := artistcontroller.ArtistControll
	ArtistController.RegisterRoute(mainRouter)

	PlayListController := playlistcontroller.PlayListControll
	PlayListController.RegisterRoute(mainRouter)

	CollectionController := collectioncontroller.CollectionControll
	CollectionController.RegisterRoute(mainRouter)

	AlbumController := albumcontroller.AlbumControll
	AlbumController.RegisterRoute(mainRouter)

	ReviewController := reviewcontroller.ReviewControll
	ReviewController.RegisterRoute(mainRouter)

	SongTypeController := songtypecontroller.SongTypeControll
	SongTypeController.RegisterRoute(mainRouter)

	CountryController := countrycontroller.CountryControll
	CountryController.RegisterRoute(mainRouter)

	OtpController := otpcontroller.OtpControll
	OtpController.RegisterRoute(mainRouter)

	// Serve static music files
	fs := http.FileServer(http.Dir("C:/Users/DPC/Desktop/MusicMp4/internal/music"))
	router.PathPrefix("/music/").Handler(http.StripPrefix("/music/", fs))

	// Start HTTP server
	http.ListenAndServe(*address, handler)
}
