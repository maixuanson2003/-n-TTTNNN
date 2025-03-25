package entity

import "time"

type Country struct {
	ID          int `gorm:"primaryKey;autoIncremen"`
	CountryName string
	CreateAt    time.Time
	Song        []Song
	Artist      []Artist
}
