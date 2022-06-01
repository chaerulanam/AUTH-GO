package repositories

import (
	"auth/models"

	"gorm.io/gorm"
)

type AuthRepo interface {
	FindAll() ([]models.User, error)
	Save(data models.User) (models.User, error)
	// UserFind(id uint) (models.User, error)
	// UserCountAllResults() ([]models.User, error)
	// UserSave(data models.User) (models.User, error)
}

func (r *userrepo) FindAll() ([]models.User, error) {
	var User []models.User

	err := r.db.Find(&User).Error

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
