package entity

import "time"

type Collection struct {
	ID             int
	NameCollection string
	CreateAt       time.Time
	UpdateAt       time.Time
	Song           []Song `gorm:"many2many:Collection_Song;constraint:OnDelete:CASCADE;"`
}
