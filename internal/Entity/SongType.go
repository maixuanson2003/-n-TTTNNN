package entity

import "time"

type SongType struct {
	ID      int `gorm:"primaryKey;autoIncremen"`
	Type    string
	CreatAt time.Time
	Song    []Song `gorm:"many2many:Song_SongType;"`
}
