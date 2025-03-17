package entity

import "time"

type Review struct {
	ID       int
	UserId   string
	SongId   int
	Content  string
	Status   string
	CreateAt time.Time
}
