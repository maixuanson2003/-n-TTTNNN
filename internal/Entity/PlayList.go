package entity

import "time"

type PlayList struct {
	ID        int
	Name      string
	CreateDay time.Time
	UpdateDay time.Time
	UserId    string
	Song      []Song `gorm:"many2many:PlayList_Song";`
}
