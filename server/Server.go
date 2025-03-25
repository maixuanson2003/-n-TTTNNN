package server

import (
	"net/http"
	database "ten_module/Database"
	albumcontroller "ten_module/internal/controller/AlbumController"
	artistcontroller "ten_module/internal/controller/ArtistController"
	Authcontroller "ten_module/internal/controller/AuthController"
	collectioncontroller "ten_module/internal/controller/CollectionController"
	historycontroller "ten_module/internal/controller/HistoryController"
	playlistcontroller "ten_module/internal/controller/PlayListController"
	reviewcontroller "ten_module/internal/controller/ReviewController"
	songcontroller "ten_module/internal/controller/SongController"
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
	database.RunSQLFile(database.Database)
	mainRouter := router.PathPrefix("/api").Subrouter()
	Cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "Accept"},
		Debug:            true,
	})
	handler := Cors.Handler(router)
	mainRouter.Use(Cors.Handler)
	// register route
	userController := UserController.UserControll
	userController.RegisterRoute(mainRouter)
	// register route
	authController := Authcontroller.AuthControllers
	authController.RegisterRoute(mainRouter)
	//song route
	songController := songcontroller.SongControllers
	songController.RegisterRoute(mainRouter)
	// history route
	HistoryController := historycontroller.HistoryControllers
	HistoryController.RegisterRoute(mainRouter)
	//artist route
	ArtistController := artistcontroller.ArtistControll
	ArtistController.RegisterRoute(mainRouter)
	//playlist route
	PlayListController := playlistcontroller.PlayListControll
	PlayListController.RegisterRoute(mainRouter)
	//collection route
	CollectionController := collectioncontroller.CollectionControll
	CollectionController.RegisterRoute(mainRouter)
	//album route
	AlbumController := albumcontroller.AlbumControll
	AlbumController.RegisterRoute(mainRouter)
	//review route
	ReviewController := reviewcontroller.ReviewControll
	ReviewController.RegisterRoute(mainRouter)
	http.ListenAndServe(*address, handler)
}
