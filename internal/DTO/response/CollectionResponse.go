package response

import "time"

type CollectionResponse struct {
	ID             int
	NameCollection string
	CreateAt       time.Time
	UpdateAt       time.Time
	Song           []SongResponse
}
