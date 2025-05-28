package request

import "time"

type AlbumUpdate struct {
	NameAlbum   string
	Description string
	ReleaseDay  time.Time
	ArtistOwner string
}
