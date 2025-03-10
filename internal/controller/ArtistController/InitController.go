package artistcontroller

import "ten_module/internal/service/artistservice"

func InitArtistControll() {
	artistservice.InitArtistService()
	InitArtistController()
}
