package response

import "time"

type AlbumResponse struct {
	ID          int
	NameAlbum   string
	Description string
	ReleaseDay  time.Time
	ArtistOwner string
	CreateDay   time.Time
	UpdateDay   time.Time
	Song        []SongResponse
	Artist      []ArtistResponse
}
