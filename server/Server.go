package server

import (
	"net/http"
	database "ten_module/Database"
	Authcontroller "ten_module/internal/controller/AuthController"
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
	mainRouter := router.PathPrefix("/api").Subrouter()
	Cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
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
	songController := songcontroller.SongControllers
	songController.RegisterRoute(mainRouter)
	http.ListenAndServe(*address, handler)
}
