package entity

import "time"

type Movie struct {
	ID            int `gorm:"primaryKey;autoIncremen"`
	NameMovie     string
	Description   string
	ReleaseDay    time.Time
	CreateDay     time.Time
	UpdateDay     time.Time
	Point         float64
	Star          float64
	Status        string
	QualityId     int
	WatchAmout    int
	MovieType     []MovieType `gorm:"many2many:Movie_MovieType;"`
	MovieResource []MovieResource
	WatchHistory  []WatchHistory
	Acting        []Acting `gorm:"many2many:Movie_Acting";`
}
