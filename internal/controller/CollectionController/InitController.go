package collectioncontroller

import collectionservice "ten_module/internal/service/CollectionService"

func InitCollectionControll() {
	collectionservice.InitCollectionService()
	InitCollectionController()
}
