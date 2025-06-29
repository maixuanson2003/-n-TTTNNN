package helper

import (
	"fmt"
	"log"
	"math"
	"ten_module/internal/repository"
)

func GetVectorFeatureForSong() (map[int][]float64, []string, error) {
	SongTypeRepo := repository.SongTypeRepo
	ArtistRepos := repository.ArtistRepo
	SongRepo := repository.SongRepo
	result := make(map[int][]float64)
	TypeArr, ErrorToGetList := SongTypeRepo.FindAll()
	SongArr, ErrorToGetSong := SongRepo.FindAll()
	ArtistArr, ErrorToGetArtist := ArtistRepos.FindAll()

	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return map[int][]float64{}, []string{}, ErrorToGetSong

	}
	if ErrorToGetArtist != nil {
		log.Print(ErrorToGetArtist)
		return map[int][]float64{}, []string{}, ErrorToGetArtist
	}
	if ErrorToGetList != nil {
		log.Print(ErrorToGetList)
		return map[int][]float64{}, []string{}, ErrorToGetList
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
		BinaryCheck := []float64{}
		SongType := SongItem.SongType
		for _, artistItem := range artist {
			check[artistItem.Name] = 1
		}
		for _, TypeItem := range SongType {
			check[TypeItem.Type] = 1
		}
		for _, Feat := range Feature {
			BinaryCheck = append(BinaryCheck, float64(check[Feat]))
		}
		result[SongItem.ID] = BinaryCheck
	}
	return result, Feature, nil
}
func GetVectorFeatureForUser(SongId int, Feature []string) ([]float64, error) {
	SongRepo := repository.SongRepo
	SongItem, ErrorToGetSong := SongRepo.GetSongById(SongId)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return []float64{}, ErrorToGetSong
	}
	Artits := SongItem.Artist
	SongType := SongItem.SongType
	BinaryCheck := []float64{}
	check := make(map[string]int)
	for _, ArtistItem := range Artits {
		check[ArtistItem.Name] = 1

	}
	for _, SongTypeItem := range SongType {
		check[SongTypeItem.Type] = 1
	}

	for _, Feat := range Feature {
		BinaryCheck = append(BinaryCheck, float64(check[Feat]))
	}

	return BinaryCheck, nil

}
func GetUserProfile(userId string, Feature []string) ([]float64, error) {
	UserRepo := repository.UserRepo
	user, errs := UserRepo.FindById(userId)
	if errs != nil {
		return []float64{}, errs
	}
	ListenHistory := user.ListenHistory
	collect := map[int]int{}
	for _, item := range ListenHistory {
		count, err := repository.ListenRepo.CountNumberSongId(item.SongId)
		if err != nil {
			return []float64{}, err
		}
		collect[item.SongId] = int(count)
	}
	userProfile := make([]float64, len(Feature))
	for Id, count := range collect {
		vector, errs := GetVectorFeatureForUser(Id, Feature)
		if errs != nil {
			return []float64{}, errs
		}
		for index, item := range vector {
			userProfile[index] += item * float64(count)
		}
	}
	log.Println("User profile vector (trước khi chuẩn hóa):")
	for i, val := range userProfile {
		log.Printf("  %s: %.4f\n", Feature[i], val)
	}
	log.Print(userProfile)
	// norm := 0.0
	// for _, val := range userProfile {
	// 	norm += val * val
	// }
	// norm = math.Sqrt(norm)
	// if norm > 0 {
	// 	for i := range userProfile {
	// 		userProfile[i] /= norm
	// 	}
	// }
	return userProfile, nil

}

func GetCosineSimilar(featureSong1 []float64, featureSong2 []float64) float64 {
	if len(featureSong1) != len(featureSong2) {
		fmt.Println("Error: Kích thước vector không khớp!")
		return 0
	}
	var dotProduct float64 = 0
	var normUser float64 = 0
	var normSong float64 = 0

	for i := 0; i < len(featureSong1); i++ {
		dotProduct += featureSong1[i] * featureSong2[i]
		normUser += math.Pow(float64(featureSong1[i]), 2)
		normSong += math.Pow(float64(featureSong2[i]), 2)
	}

	normUser = math.Sqrt(normUser)
	normSong = math.Sqrt(normSong)
	// s := float64(dotProduct) / (normUser * normSong)

	if normUser == 0 || normSong == 0 {
		return 0
	}
	return float64(dotProduct) / (normUser * normSong)
}
