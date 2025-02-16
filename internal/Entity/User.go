package entity

type User struct {
	ID           string `gorm:"type:varchar(255);primaryKey";`
	Username     string
	FullName     string
	Password     string
	Phone        string
	Email        string
	Address      string
	Gender       string
	Age          string
	Role         string
	WatchHistory []WatchHistory
}
