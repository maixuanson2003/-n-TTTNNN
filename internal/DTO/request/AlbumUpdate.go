package request

import "time"

type AlbumUpdate struct {
	NameAlbum   string
	Description string
	ReleaseDay  time.Time
	ArtistOwner string
	Artist      []int `validate:"required,min=1,dive,gt=0"`
}
