package repositories

import (
	"auth/models"

	"gorm.io/gorm"
)

type AuthRepo interface {
	FindAll() ([]models.User, error)
	IsRegistered(email string, username string) (models.User, error)
	Save(data models.User) (models.User, error)
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

type userrepo struct {
	db *gorm.DB
}

func UserRepository(db *gorm.DB) *userrepo {
	return &userrepo{db}
}
