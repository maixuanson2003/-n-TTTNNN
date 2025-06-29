package entity

import "time"

type Album struct {
	ID          int
	NameAlbum   string
	Image       string
	Description string
	ReleaseDay  time.Time
	ArtistOwner string
	CreateDay   time.Time
	UpdateDay   time.Time
	Song        []Song
	Artist      []Artist `gorm:"many2many:Album_Artist;constraint:OnDelete:CASCADE;"`
}
