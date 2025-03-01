package songcontroller

import "ten_module/internal/service/songservice"

func InitSongService() {
	songservice.InitSongService()
	InitSongController()
}
