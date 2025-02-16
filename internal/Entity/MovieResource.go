package entity

import "time"

type MovieResource struct {
	ID      int `gorm:"primaryKey;autoIncremen"`
	Chapter string
	Video   string
	CreatAt time.Time
	MovieId int
}
