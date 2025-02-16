package entity

import "time"

type MovieType struct {
	ID      int `gorm:"primaryKey;autoIncremen"`
	Type    string
	CreatAt time.Time
	Movie   []Movie `gorm:"many2many:Movie_MovieType;"`
}
