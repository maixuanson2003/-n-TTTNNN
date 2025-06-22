package reviewservice

import (
	"errors"
	"log"
	"ten_module/internal/DTO/request"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/repository"
	"time"
)

type ReviewService struct {
	ReviewRepo *repository.ReviewRepository
	SongRepo   *repository.SongRepository
	UserRepo   *repository.UserRepository
}
type MessageResponse struct {
	Message string
	Status  string
}

const (
	Status_Up   = "PUBLISH"
	Status_Down = "NOT_PUBLISH"
)

var ReviewServe *ReviewService

func InitReviewService() {
	ReviewServe = &ReviewService{
		ReviewRepo: repository.ReviewRepo,
		SongRepo:   repository.SongRepo,
		UserRepo:   repository.UserRepo,
	}
}

type ReviewServiceInterface interface {
	GetListReview() ([]response.ReviewResponse, error)
	GetListReviewBySong(SongId int) ([]response.ReviewResponse, error)
	CreateReview(Review request.ReviewRequest) (MessageResponse, error)
	UpdateReview(Review request.ReviewRequest, ReviewId int) (MessageResponse, error)
}

func MapReviewEntityToReviewResponse(Review entity.Review, UserName string) response.ReviewResponse {
	return response.ReviewResponse{
		Id:       Review.ID,
		UserName: UserName,
		Content:  Review.Content,
		Status:   Review.Status,
		CreateAt: Review.CreateAt,
	}

}
func MapReviewRequestToReviewEntity(Review request.ReviewRequest) entity.Review {
	return entity.Review{
		UserId:   Review.UserId,
		SongId:   Review.SongId,
		Content:  Review.Content,
		CreateAt: time.Now(),
		Status:   Status_Up,
	}

}
func (ReviewServe *ReviewService) GetListReview() ([]response.ReviewResponse, error) {
	ReviewRepo := ReviewServe.ReviewRepo
	UserRepo := ReviewServe.UserRepo
	ReviewList, ErrorToGetListReview := ReviewRepo.FindAll()

	if ErrorToGetListReview != nil {
		log.Print(ErrorToGetListReview)
		return nil, ErrorToGetListReview
	}
	ReviewResponse := []response.ReviewResponse{}
	for _, ReviewItem := range ReviewList {
		User, ErrorToGetUser := UserRepo.FindById(ReviewItem.UserId)
		if ErrorToGetUser != nil {
			log.Print(ErrorToGetUser)
			return nil, ErrorToGetUser
		}
		ReviewResponse = append(ReviewResponse, MapReviewEntityToReviewResponse(ReviewItem, User.Username))

	}
	return ReviewResponse, nil

}
func (ReviewServe *ReviewService) GetListReviewBySong(SongId int) ([]response.ReviewResponse, error) {
	SongRepo := ReviewServe.SongRepo
	UserRepo := ReviewServe.UserRepo
	Song, ErrorToGetSong := SongRepo.GetSongById(SongId)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return nil, ErrorToGetSong
	}
	ReviewArray := Song.Review
	ReviewResponse := []response.ReviewResponse{}
	for _, ReviewItem := range ReviewArray {
		if ReviewItem.Status == Status_Up {
			User, ErrorToGetUser := UserRepo.FindById(ReviewItem.UserId)
			if ErrorToGetUser != nil {
				log.Print(ErrorToGetUser)
				return nil, ErrorToGetUser
			}
			ReviewResponse = append(ReviewResponse, MapReviewEntityToReviewResponse(ReviewItem, User.Username))
		}

	}
	return ReviewResponse, nil
}
func (ReviewServe *ReviewService) CreateReview(Review request.ReviewRequest) (MessageResponse, error) {
	ReviewRepo := ReviewServe.ReviewRepo
	ReviewEntity := MapReviewRequestToReviewEntity(Review)
	ErrorToCreateReview := ReviewRepo.CreateReview(ReviewEntity)
	if ErrorToCreateReview != nil {
		log.Print(ErrorToCreateReview)
		return MessageResponse{
			Message: "failed to create review",
			Status:  "Failed",
		}, ErrorToCreateReview
	}
	return MessageResponse{
		Message: "success to create review",
		Status:  "Success",
	}, nil

}
func mapStatus(input string) string {
	switch input {
	case "PUBLISH":
		return Status_Up
	case "NOT_PUBLISH":
		return Status_Down
	default:
		return Status_Down
	}
}
func (ReviewServe *ReviewService) UpdateReview(Status string, ReviewId int) (MessageResponse, error) {
	ReviewRepo := ReviewServe.ReviewRepo
	log.Print(Status)
	Review, errors := ReviewRepo.FindById(ReviewId)
	if errors != nil {
		return MessageResponse{
			Message: "failed",
			Status:  "Failed",
		}, errors
	}
	log.Print(mapStatus(Status))
	Review.Status = mapStatus(Status)
	ErrorToUpdateReview := ReviewRepo.UpdateReview(Review, ReviewId)
	if ErrorToUpdateReview != nil {
		log.Print(ErrorToUpdateReview)
		return MessageResponse{
			Message: "failed to update review",
			Status:  "Failed",
		}, ErrorToUpdateReview
	}
	return MessageResponse{
		Message: "success to update review",
		Status:  "Success",
	}, nil

}
func (ReviewServe *ReviewService) DeleteReviewById(ReviewId int) (MessageResponse, error) {
	ReviewRepo := ReviewServe.ReviewRepo
	errors := ReviewRepo.DeleteReview(ReviewId)
	if errors != nil {
		log.Print(errors)
		return MessageResponse{}, errors
	}
	return MessageResponse{
		Message: "failed",
		Status:  "Failed",
	}, nil

}
func (ReviewServe *ReviewService) DeleteReviewByIdForUser(reviewId int, userId string) (MessageResponse, error) {
	reviewRepo := ReviewServe.ReviewRepo
	userRepo := ReviewServe.UserRepo

	user, err := userRepo.FindById(userId)
	if err != nil {
		return MessageResponse{}, err
	}

	// Kiểm tra xem reviewId có thuộc user không
	owned := false
	for _, review := range user.Review {
		if review.ID == reviewId {
			owned = true
			break
		}
	}

	if !owned {
		return MessageResponse{}, errors.New("review does not belong to this user")
	}
	err = reviewRepo.DeleteReview(reviewId)
	if err != nil {
		log.Print(err)
		return MessageResponse{}, err
	}
	return MessageResponse{
		Message: "Review deleted successfully",
		Status:  "Success",
	}, nil

}
func (ReviewServe *ReviewService) GetListReviewForUser(userId string) ([]response.ReviewResponse, error) {
	UserRepo := ReviewServe.UserRepo
	User, ErrorToGetSong := UserRepo.FindById(userId)
	if ErrorToGetSong != nil {
		log.Print(ErrorToGetSong)
		return nil, ErrorToGetSong
	}
	ReviewArray := User.Review
	ReviewResponse := []response.ReviewResponse{}
	for _, ReviewItem := range ReviewArray {
		if ReviewItem.Status == Status_Up {
			User, ErrorToGetUser := UserRepo.FindById(ReviewItem.UserId)
			if ErrorToGetUser != nil {
				log.Print(ErrorToGetUser)
				return nil, ErrorToGetUser
			}
			ReviewResponse = append(ReviewResponse, MapReviewEntityToReviewResponse(ReviewItem, User.Username))
		}

	}
	return ReviewResponse, nil
}
