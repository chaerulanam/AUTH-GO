package services

import (
	"auth/models"
	"auth/repositories"
	"auth/requests"
	"database/sql"
)

type AuthService interface {
	FindAll() ([]models.User, error)
	Save(data requests.AuthReq) (models.User, error)
}

func (s *authservice) FindAll() ([]models.User, error) {
	user, err := s.userrepository.FindAll()
	return user, err
}

func (s *authservice) Save(data requests.AuthReq) (models.User, error) {

	// passwordHash, _ := helper.HashPassword(data.Password)
	var activatehas = sql.NullString{}

	// if config.AUTH_EMAIL_VERIFIKASI == true {
	// 	activatehas = sql.NullString{String: helper.GenerateToken(data.Email), Valid: true}
	// }

	userModel := models.User{
		// Username:       data.Username,
		// Email:          data.Email,
		// PasswordHash:   passwordHash,
		ResetAt:        sql.NullTime{},
		ResetExpire:    sql.NullTime{},
		ActivateHash:   activatehas,
		Status:         sql.NullString{},
		StatusMessage:  sql.NullString{},
		Active:         true,
		ForcePassReset: false,
	}

	user, err := s.userrepository.Save(userModel)
	return user, err
}

type authservice struct {
	userrepository repositories.AuthRepo
}

func AuthServ(userrepository repositories.AuthRepo) *authservice {

	return &authservice{userrepository}
}
