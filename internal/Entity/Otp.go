package entity

import "time"

type Otp struct {
	ID        int `gorm:"primaryKey;autoIncremen"`
	Otp       string
	Create_at time.Time
}
