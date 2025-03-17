package response

import "time"

type ReviewResponse struct {
	Id       int
	UserName string
	Content  string
	Status   string
	CreateAt time.Time
}
