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
	SongType      []SongType `gorm:"many2many:Song_SongType;"`
	SongResource  string     `gorm:"not null"`
	ListenHistory []ListenHistory
	Artist        []Artist     `gorm:"many2many:Song_Artist;"`
	User          []User       `gorm:"many2many:User_Like;"`
	PlayList      []PlayList   `gorm:"many2many:PlayList_Song;"`
	Collection    []Collection `gorm:"many2many:Collection_Song;"`
}
