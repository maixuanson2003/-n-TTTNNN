package entity

type Acting struct {
	ID          int
	Name        string
	BirthDay    string
	Description string
	Country     string
	Movie       []Movie `gorm:"many2many:Movie_Acting";`
}
