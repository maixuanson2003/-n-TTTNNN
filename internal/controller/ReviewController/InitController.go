package reviewcontroller

import "ten_module/internal/service/reviewservice"

func InitReviewControll() {
	reviewservice.InitReviewService()
	InitReviewController()
}
