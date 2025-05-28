package request

import (
	"mime/multipart"
	"time"
)

type AlbumRequest struct {
	NameAlbum   string        `validate:"required,min=2,max=100"`
	Description string        `validate:"max=500"`
	ReleaseDay  time.Time     `validate:"required"`
	ArtistOwner string        `validate:"required"`
	Song        []SongRequest `validate:"required,dive"`
	Artist      []int         `validate:"required,min=1,dive,gt=0"`
}
type SongFileAlbum struct {
	SongName string
	File     multipart.File
}
