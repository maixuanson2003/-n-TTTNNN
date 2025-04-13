package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
func RunSQLFile(db *gorm.DB) error {
	FolderPath := "C:\\Users\\DPC\\Desktop\\MusicMp4\\Database"
	fileArry, errorToGetFileArray := os.ReadDir("C:\\Users\\DPC\\Desktop\\MusicMp4\\Database")
	if errorToGetFileArray != nil {
		return errorToGetFileArray
	}
	for _, fileItem := range fileArry {
		filePath := filepath.Join(FolderPath, fileItem.Name())
		fileext := filepath.Ext(filePath)
		if fileext == ".sql" {
			sqlBytes, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			sqlContent := string(sqlBytes)
			err = db.Exec(sqlContent).Error
			if err != nil {
				return err
			}
		}

	}
	fmt.Println("✅ Seed dữ liệu thành công từ file:", "ok")
	return nil
}
