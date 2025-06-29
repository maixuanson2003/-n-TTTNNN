package response

import "time"

type AlbumResponse struct {
	ID          int
	NameAlbum   string
	Description string
	ReleaseDay  time.Time
	ArtistOwner string
	CreateDay   time.Time
	Image       string
	UpdateDay   time.Time
	Song        []SongResponseAlbum
	Artist      []ArtistResponse
}
