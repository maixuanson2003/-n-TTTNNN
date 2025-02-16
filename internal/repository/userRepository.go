package repository

import (
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

type Interface interface {
	FindAll() ([]entity.User, error)
	FindById(Id int) (entity.User, error)
	Create(User entity.User) error
	Update(User entity.User) error
	DeleteById(Id int) error
	DeleteAll(User []entity.User) error
}

func (instance *UserRepository) FindAll() ([]entity.User, error) {
	Database := instance.DB
	var User []entity.User
	err := Database.Model(&entity.User{}).Find(&User).Error
	if err != nil {
		return nil, err
	}
	return User, nil
}
func (instance *UserRepository) FindById(Id int) (entity.User, error) {
	Database := instance.DB
	var User entity.User
	err := Database.First(&User, Id).Error
	if err != nil {
		return entity.User{}, err
	}
	return User, nil
}
func (instance *UserRepository) Create(User entity.User) error {
	Database := instance.DB
	err := Database.Create(&User).Error
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
func (instance *UserRepository) Update(User entity.User) error {
	Database := instance.DB
	err := Database.Save(&User).Error
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
func (instance *UserRepository) DeleteById(Id int) error {
	Database := instance.DB
	err := Database.Delete(&entity.User{}, Id).Error
	if err != nil {
		log.Print(err)
		return err
	}
	return nil

}
func (instance *UserRepository) DeleteAll(User []entity.User) error {
	Database := instance.DB
	err := Database.Delete(&User).Error
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
