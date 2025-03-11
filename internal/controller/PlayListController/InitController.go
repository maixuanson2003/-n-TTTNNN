package playlistcontroller

import "ten_module/internal/service/playlistservice"

func InitPlayListControll() {
	playlistservice.InitPlayListService()
	InitPlayListController()

}
