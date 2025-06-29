package repository

import (
	"fmt"
	"log"
	database "ten_module/Database"
	entity "ten_module/internal/Entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

var UserRepo *UserRepository

func InitUserRepo() {
	UserRepo = &UserRepository{
		DB: database.Database,
	}
}

type UserRepoInterface interface {
	FindAll() ([]entity.User, error)
	FindById(Id int) (entity.User, error)
	Create(User entity.User) error
	Update(User entity.User, id string) error
	DeleteById(Id int) error
	DeleteAll(User []entity.User) error
	GetUserQuery(Name string, Age int, Email string, Address string, Role string, Gender string) ([]entity.User, error)
}

func (instance *UserRepository) FindAll() ([]entity.User, error) {
	Database := instance.DB
	var User []entity.User
	err := Database.Model(&entity.User{}).Order("id DESC").Find(&User).Error
	if err != nil {
		return nil, err
	}
	return User, nil
}
func (instance *UserRepository) FindById(Id string) (entity.User, error) {
	Database := instance.DB
	var User entity.User
	err := Database.
		Preload("ListenHistory").
		Preload("Song").
		Preload("PlayList").
		Preload("Review").
		Where("id = ?", Id).
		First(&User).Error
	if err != nil {
		return entity.User{}, err
	}
	return User, nil
}
func (instance *UserRepository) FindByEmail(email string) (entity.User, error) {
	Database := instance.DB
	var user entity.User
	err := Database.
		Preload("ListenHistory").
		Preload("Song").
		Preload("PlayList").
		Preload("Review").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
func (instance *UserRepository) Create(User entity.User) error {
	Database := instance.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Create(&User).Error
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
func (instance *UserRepository) Update(User entity.User, id string) error {
	Database := instance.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Where("id=?", id).Save(&User).Error
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
func (instance *UserRepository) DeleteById(Id string) error {
	db := instance.DB

	var user entity.User
	errFind := db.Preload("ListenHistory").
		Preload("Song").
		Preload("PlayList").
		Preload("Review").
		Where("id = ?", Id).
		First(&user).Error

	if errFind != nil {
		return errFind
	}

	errDelete := db.Select("ListenHistory", "Song", "PlayList", "Review").Delete(&user).Error
	if errDelete != nil {
		return errDelete
	}

	return nil

}
func (instance *UserRepository) DeleteAll(User []entity.User) error {
	Database := instance.DB
	errs := Database.Transaction(func(tx *gorm.DB) error {
		err := Database.Delete(&User).Error
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
func (instance *UserRepository) GetUserQuery(Name string, Age int, Email string, Address string, Role string, Gender string) ([]entity.User, error) {
	Database := instance.DB
	var User []entity.User
	Query := map[string]interface{}{}
	if Name != "" {
		Query["full_name"] = Name

	}
	if Age > 0 {
		Query["age"] = Age
	}
	if Email != "" {
		Query["email"] = Email
	}
	if Address != "" {
		Query["address"] = Address
	}
	if Role != "" {
		Query["role"] = Role
	}
	if Gender != "" {
		Query["gender"] = Gender
	}
	fmt.Print(Query)
	err := Database.Where(Query).Find(&User).Error
	if err != nil {
		return nil, err
	}
	return User, nil
}
func (instance *UserRepository) DeleteSongLike(userID string, songID int) error {
	result := instance.DB.
		Table("user_likes").
		Where("user_id = ? AND song_id = ?", userID, songID).
		Delete(nil)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no like found for user_id=%s and song_id=%d", userID, songID)
	}
	return nil
}
