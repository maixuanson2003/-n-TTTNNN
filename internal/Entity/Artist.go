package entity

type Artist struct {
	ID          int
	Name        string
	Image       string
	BirthDay    string
	Description string
	CountryId   int
	Song        []Song  `gorm:"many2many:Song_Artist;constraint:OnDelete:CASCADE;"`
	Album       []Album `gorm:"many2many:Album_Artist;constraint:OnDelete:CASCADE;"`
}
