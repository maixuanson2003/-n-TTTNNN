package entity

type Artist struct {
	ID          int
	Name        string
	BirthDay    string
	Description string
	CountryId   int
	Song        []Song  `gorm:"many2many:Song_Artist;"`
	Album       []Album `gorm:"many2many:Album_Artist;"`
}
