package helper

import (
	"fmt"
	"log"
	"math"
	"ten_module/internal/repository"
)

func GetVectorFeatureForSong() (map[int][]int16, []string, error) {
	SongTypeRepo := repository.SongTypeRepo
	ArtistRepos := repository.ArtistRepo
	SongRepo := repository.SongRepo
	result := make(map[int][]int16)
	TypeArr, ErrorToGetList := SongTypeRepo.FindAll()
	SongArr, ErrorToGetSong := SongRepo.FindAll()
	ArtistArr, ErrorToGetArtist := ArtistRepos.FindAll()

	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return map[int][]int16{}, []string{}, ErrorToGetSong

	}
	if ErrorToGetArtist != nil {
		log.Print(ErrorToGetArtist)
		return map[int][]int16{}, []string{}, ErrorToGetArtist
	}
	if ErrorToGetList != nil {
		log.Print(ErrorToGetList)
		return map[int][]int16{}, []string{}, ErrorToGetList
	}
	Feature := []string{}
	for _, value := range ArtistArr {
		Feature = append(Feature, value.Name)
	}
	for _, value := range TypeArr {
		Feature = append(Feature, value.Type)
	}
	for _, SongItem := range SongArr {
		check := make(map[string]int)
		artist := SongItem.Artist
		BinaryCheck := []int16{}
		SongType := SongItem.SongType
		for _, artistItem := range artist {
			check[artistItem.Name] = 1
		}
		for _, TypeItem := range SongType {
			check[TypeItem.Type] = 1
		}
		for _, Feat := range Feature {
			BinaryCheck = append(BinaryCheck, int16(check[Feat]))
		}
		result[SongItem.ID] = BinaryCheck
	}
	return result, Feature, nil
}
func GetVectorFeatureForUser(SongId int, Feature []string) ([]int16, error) {
	SongRepo := repository.SongRepo
	SongItem, ErrorToGetSong := SongRepo.GetSongById(SongId)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return []int16{}, ErrorToGetSong
	}
	Artits := SongItem.Artist
	SongType := SongItem.SongType
	BinaryCheck := []int16{}
	check := make(map[string]int)
	for _, ArtistItem := range Artits {
		check[ArtistItem.Name] = 1

	}
	for _, SongTypeItem := range SongType {
		check[SongTypeItem.Type] = 1
	}

	for _, Feat := range Feature {
		BinaryCheck = append(BinaryCheck, int16(check[Feat]))
	}
	fmt.Print(BinaryCheck)
	return BinaryCheck, nil

}

func GetCosineSimilar(featureSong1 []int16, featureSong2 []int16) float64 {
	if len(featureSong1) != len(featureSong2) {
		fmt.Println("Error: Kích thước vector không khớp!")
		return 0
	}
	var dotProduct int16 = 0
	var normUser float64 = 0
	var normSong float64 = 0

	for i := 0; i < len(featureSong1); i++ {
		dotProduct += featureSong1[i] * featureSong2[i]
		normUser += math.Pow(float64(featureSong1[i]), 2)
		normSong += math.Pow(float64(featureSong2[i]), 2)
	}

	normUser = math.Sqrt(normUser)
	normSong = math.Sqrt(normSong)
	s := float64(dotProduct) / (normUser * normSong)
	fmt.Println(s)

	if normUser == 0 || normSong == 0 {
		return 0
	}
	return float64(dotProduct) / (normUser * normSong)
}
