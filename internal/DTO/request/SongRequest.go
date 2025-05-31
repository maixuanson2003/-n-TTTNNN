package request

import (
	"mime/multipart"
	"time"
)

type SongRequest struct {
	NameSong    string    `validate:"required,min=2,max=100"`
	Description string    `validate:"max=500"`
	ReleaseDay  time.Time `validate:"required"`
	Point       float64   `validate:"gte=0,lte=10"`
	Status      string
	CountryId   int   `validate:"required"`
	SongType    []int `validate:"required,min=1,dive,gt=0"`
	Artist      []int `validate:"required,min=1,dive,gt=0"`
}
type SongFile struct {
	File multipart.File
}
