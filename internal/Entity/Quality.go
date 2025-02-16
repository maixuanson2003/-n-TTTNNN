package entity

import "time"

type Quality struct {
	ID          int `gorm:"primaryKey;autoIncremen"`
	QualityName string
	CreateAt    time.Time
	Movie       []Movie
}
