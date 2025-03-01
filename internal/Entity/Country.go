package entity

import "time"

type Country struct {
	ID          int `gorm:"primaryKey;autoIncremen"`
	QualityName string
	CreateAt    time.Time
	Song        []Song
}
