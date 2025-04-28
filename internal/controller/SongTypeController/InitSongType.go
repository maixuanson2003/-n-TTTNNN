package songtypecontroller

import "ten_module/internal/service/songtypeservice"

func InitSongTypeControll() {
	songtypeservice.InitSongTypeService()
	InitSongTypeController()
}
