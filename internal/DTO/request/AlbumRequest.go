package request

import (
	"mime/multipart"
	"time"
)

type AlbumRequest struct {
	NameAlbum   string
	Description string
	ReleaseDay  time.Time
	ArtistOwner string
	Song        []SongRequest
	Artist      []int
}
type SongFileAlbum struct {
	SongName string
	File     multipart.File
}
