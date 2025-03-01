package database

import (
	"log"
	entity "ten_module/internal/Entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Init() {
	DatabaseInit := "root:Maixuanson2003@tcp(127.0.0.1:3306)/musicdb?charset=utf8mb4&parseTime=True&loc=Local"
	appDB, err := gorm.Open(mysql.Open(DatabaseInit), &gorm.Config{})
	if err != nil {
		log.Print(err)
		return
	}
	Database = appDB
}

func MigrateDB(DB *gorm.DB) {
	tables := []interface{}{
		&entity.Album{},
		&entity.PlayList{},
		&entity.Review{},
		&entity.Collection{},
		&entity.User{},
		&entity.Song{},
		&entity.SongType{},
		&entity.Artist{},
		&entity.Country{},
		&entity.ListenHistory{},
	}
	for _, table := range tables {
		DB.AutoMigrate(table)
	}
}
