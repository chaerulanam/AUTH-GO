package repositories

import (
	"auth/models"

	"gorm.io/gorm"
)

type AuthRepo interface {
	FindAll() ([]models.User, error)
	IsRegistered(email string, username string) (models.User, error)
	Save(data models.User) (models.User, error)
	SaveAuthLogin(data models.AuthLogin) (models.AuthLogin, error)
	AddGroup(data models.AuthGroup) (models.AuthGroup, error)
	AddPermission(data models.AuthPermission) (models.AuthPermission, error)
}

func (r *userrepo) FindAll() ([]models.User, error) {
	var User []models.User

	err := r.db.Find(&User).Error

	return User, err
}

func (r *userrepo) IsRegistered(email string, username string) (models.User, error) {

	var User models.User

	err := r.db.Where("email = ?", email).Or("username = ?", username).Find(&User).Error

	return User, err

}

func (r *userrepo) Save(data models.User) (models.User, error) {
	err := r.db.Create(&data).Error
	return data, err
}

func (r *userrepo) SaveAuthLogin(data models.AuthLogin) (models.AuthLogin, error) {
	err := r.db.Create(&data).Error
	return data, err
}

func (r *userrepo) AddGroup(data models.AuthGroup) (models.AuthGroup, error) {
	err := r.db.Create(&data).Error
	return data, err
}

func (r *userrepo) AddPermission(data models.AuthPermission) (models.AuthPermission, error) {
	err := r.db.Create(&data).Error
	return data, err
}

type userrepo struct {
	db *gorm.DB
}

func UserRepository(db *gorm.DB) *userrepo {
	return &userrepo{db}
}
