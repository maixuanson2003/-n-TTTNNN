package request

import (
	"mime/multipart"
	"reflect"
	"time"
)

type SongRequest struct {
	NameSong    string    `validate:"required,min=2,max=100"`
	Description string    `validate:"max=500"`
	ReleaseDay  time.Time `validate:"required"`
	Point       float64   `validate:"gte=0,lte=5"`
	Status      string
	CountryId   int   `validate:"required"`
	SongType    []int `validate:"required,min=1,dive,gt=0"`
	Artist      []int `validate:"required,min=1,dive,gt=0"`
}
type SongRequestUpdate struct {
	ID          *int
	NameSong    string    `validate:"required,min=2,max=100"`
	Description string    `validate:"max=500"`
	ReleaseDay  time.Time `validate:"required"`
	Point       float64   `validate:"gte=0,lte=5"`
	Status      string
	CountryId   int   `validate:"required"`
	SongType    []int `validate:"required,min=1,dive,gt=0"`
	Artist      []int `validate:"required,min=1,dive,gt=0"`
}
type SongFile struct {
	File *multipart.FileHeader
}

func IsTypedNil(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	}
	return false
}
