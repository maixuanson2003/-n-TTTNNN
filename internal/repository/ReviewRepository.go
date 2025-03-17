package repository

import (
	"log"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	DB *gorm.DB
}

var ReviewRepo *ReviewRepository

func InitReviewRepository() {
	ReviewRepo = &ReviewRepository{
		DB: database.Database,
	}
}

type ReviewRepositoryInterface interface {
	FindAll() ([]entity.Review, error)
	GetAlbumById(Id int) (entity.Album, error)
	CreateReview(Review entity.Review) error
	UpdateReview(Review entity.Review, id int) error
	DeleteReviewById(Id int) error
}

func (ReviewRepo *ReviewRepository) FindAll() ([]entity.Review, error) {
	Database := ReviewRepo.DB
	var Review []entity.Review
	err := Database.Model(&entity.Review{}).Find(&Review).Error
	if err != nil {
		return nil, err
	}
	return Review, nil
}
func (ReviewRepo *ReviewRepository) CreateReview(Review entity.Review) error {
	Database := ReviewRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Create(&Review).Error
		if err != nil {
			return err

		}
		return nil

	})
	if errs != nil {
		log.Print(errs)
		return errs
	}
	return nil
}
func (ReviewRepo *ReviewRepository) UpdateReview(Review entity.Review, id int) error {
	Database := ReviewRepo.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Where("id=?", id).Save(&Review).Error
		if err != nil {
			return err
		}
		return nil
	})
	if errs != nil {
		log.Print(errs)
		return errs
	}
	return nil
}
