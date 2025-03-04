package historycontroller

import (
	"ten_module/internal/service/listenhistoryservice"
)

func InitHistoryService() {
	listenhistoryservice.InitListenHistoryService()
	InitHistoryControllers()
}
