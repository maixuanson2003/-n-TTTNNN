package entity

import "time"

type Review struct {
	ID       int
	UserId   int
	Content  string
	CreateAt time.Time
}
