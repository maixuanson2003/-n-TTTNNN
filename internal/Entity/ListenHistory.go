package entity

import "time"

type ListenHistory struct {
	ID     int `gorm:"primaryKey;autoIncrement"`
	SongId int
	UserId *string `gorm:"type:varchar(255);index"`

	ListenDay time.Time
}
