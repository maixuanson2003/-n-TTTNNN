package request

import (
	"mime/multipart"
	"time"
)

type SongRequest struct {
	NameSong    string
	Description string
	ReleaseDay  time.Time
	Point       float64
	Status      string
	CountryId   int
	SongType    []int
	Artist      []int
}
type SongFile struct {
	File multipart.File
}
