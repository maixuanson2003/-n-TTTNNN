package albumcontroller

import "ten_module/internal/service/albumservice"

func InitAlbumControll() {
	albumservice.InitAlbumSerivce()
	InitAlbumController()
}
