package response

import "time"

type SongResponse struct {
	ID           int
	NameSong     string
	Description  string
	ReleaseDay   time.Time
	CreateDay    time.Time
	UpdateDay    time.Time
	Point        float64
	LikeAmount   int
	Status       string
	CountryId    int
	ListenAmout  int
	AlbumId      int
	SongResource string
}
