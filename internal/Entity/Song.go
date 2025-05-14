package entity

import "time"

type Song struct {
	ID            int `gorm:"primaryKey;autoIncremen"`
	NameSong      string
	Description   string
	ReleaseDay    time.Time
	CreateDay     time.Time
	UpdateDay     time.Time
	Point         float64
	LikeAmount    int
	Status        string
	CountryId     int
	ListenAmout   int
	AlbumId       *int       `gorm:"default:null"`
	SongType      []SongType `gorm:"many2many:Song_SongType;constraint:OnDelete:CASCADE;"`
	SongResource  string     `gorm:"not null"`
	ListenHistory []ListenHistory
	Review        []Review
	Artist        []Artist     `gorm:"many2many:Song_Artist;constraint:OnDelete:CASCADE;"`
	User          []User       `gorm:"many2many:User_Like;constraint:OnDelete:CASCADE;"`
	PlayList      []PlayList   `gorm:"many2many:PlayList_Song;constraint:OnDelete:CASCADE;"`
	Collection    []Collection `gorm:"many2many:Collection_Song;constraint:OnDelete:CASCADE;"`
}
