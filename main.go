package main

import (
	database "ten_module/Database"
	"ten_module/server"
)

func main() {
	server.InitSingleton()

	server := server.GetNewServer(":8080", database.Database)
	server.Run(&server.Address, server.Database)
}
