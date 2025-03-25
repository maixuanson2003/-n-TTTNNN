package response

import "time"

type CountryResponse struct {
	Id          int
	CountryName string
	CreateAt    time.Time
}
