package server

import (
	"net/http"
	database "ten_module/Database"
	Authcontroller "ten_module/internal/controller/AuthController"
	"ten_module/internal/controller/UserController"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

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
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	mainRouter := router.PathPrefix("/api").Subrouter()
	// register route
	userController := UserController.UserControll
	userController.RegisterRoute(mainRouter)
	// register route
	authController := Authcontroller.AuthControllers
	authController.RegisterRoute(mainRouter)
	http.ListenAndServe(*address, mainRouter)
}
