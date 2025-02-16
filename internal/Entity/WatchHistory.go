package entity

import "time"

type WatchHistory struct {
	ID      int `gorm:"primaryKey;autoIncrement"`
	MovieId int
	UserId  string `gorm:"type:varchar(255);index"`

	WatchDay time.Time
}
