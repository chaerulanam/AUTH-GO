package repositories

import (
	"auth/V1/Auth/dto"
	"auth/models"

	"gorm.io/gorm"
)

type AuthRepo interface {
	FindAll() ([]models.User, error)
	IsRegistered(data models.User) (models.User, error)
	Save(data models.User) (models.User, error)
	SaveAuthLogin(data models.AuthLogin) (models.AuthLogin, error)
	AddGroup(data models.AuthGroup) (models.AuthGroup, error)
	AddPermission(data models.AuthPermission) (models.AuthPermission, error)
	GetGroupId(data string) (models.AuthGroup, error)
	FindGroupId(name string) (models.AuthGroup, error)
	AddUserToGroup(data models.AuthGroupUser) (models.AuthGroupUser, error)
	RemoveUserFromGroup(data models.AuthGroupUser) (models.AuthGroupUser, error)
	CountUsers() (int64, error)
	CountUserBySearch(data dto.DatatablesReq) (int64, error)
	DatatablesFind(data dto.DatatablesReq) ([]models.User, error)
	DatatablesSearch(data dto.DatatablesReq) ([]models.User, error)
}

func (r *userrepo) FindAll() ([]models.User, error) {
	var User []models.User

	err := r.db.Preload("AuthGroupUser.AuthGroup").Find(&User).Error
	return User, err
}

func (r *userrepo) IsRegistered(data models.User) (models.User, error) {

	var User models.User

	err := r.db.Where("email = ?", data.Email).Or("username = ?", data.Username).Find(&User).Error

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

func (r *userrepo) GetGroupId(data string) (models.AuthGroup, error) {

	var _data models.AuthGroup

	err := r.db.Where("name = ?", data).Find(&_data).Error

	return _data, err
}

func (r *userrepo) FindGroupId(name string) (models.AuthGroup, error) {
	var data models.AuthGroup

	err := r.db.Where("name = ?", name).Find(&data).Error
	return data, err
}

func (r *userrepo) AddUserToGroup(data models.AuthGroupUser) (models.AuthGroupUser, error) {
	err := r.db.Create(&data).Error
	return data, err
}

func (r *userrepo) RemoveUserFromGroup(data models.AuthGroupUser) (models.AuthGroupUser, error) {
	err := r.db.Where("user_id = ?", data.UserID).Where("group_id = ?", data.GroupID).Delete(&data).Error
	return data, err
}

func (r *userrepo) CountUsers() (int64, error) {
	var count int64
	var User []models.User

	err := r.db.Model(&User).Count(&count).Error
	return count, err
}

func (r *userrepo) CountUserBySearch(data dto.DatatablesReq) (int64, error) {
	var count int64
	var User []models.User

	err := r.db.Model(&User).Where("Username LIKE ?", "%"+data.SearchValue+"%").Or("Email LIKE ?", "%"+data.SearchValue+"%").Count(&count).Error
	return count, err
}

func (r *userrepo) DatatablesFind(data dto.DatatablesReq) ([]models.User, error) {
	var User []models.User
	var column = [4]string{"id", "username", "email"}

	err := r.db.Preload("AuthGroupUser.AuthGroup").Order(column[data.Order0Column] + " " + data.Order0Dir).Limit(data.Length).Offset(data.Start).Find(&User).Error
	return User, err
}

func (r *userrepo) DatatablesSearch(data dto.DatatablesReq) ([]models.User, error) {
	var User []models.User
	var column = [4]string{"id", "username", "email"}

	err := r.db.Preload("AuthGroupUser.AuthGroup").Where("Username LIKE ?", "%"+data.SearchValue+"%").Or("Email LIKE ?", "%"+data.SearchValue+"%").Order(column[data.Order0Column] + " " + data.Order0Dir).Limit(data.Length).Offset(data.Start).Find(&User).Error
	return User, err
}

type userrepo struct {
	db *gorm.DB
}

func UserRepository(db *gorm.DB) *userrepo {
	return &userrepo{db}
}
