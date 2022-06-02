package services

import (
	"auth/config"
	"auth/helper"
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

	passwordHash, _ := helper.HashPassword(data.Password)
	var activatehas = sql.NullString{}

	if config.AUTH_EMAIL_VERIFIKASI == true {
		activatehas = sql.NullString{String: helper.GenerateToken(data.Email), Valid: true}
		var htmlaktivasi = "<p>Ini adalah email aktivasi akun anda untuk website " + config.APP_NAME +
			".</p><p>Untuk mengaktifkan akun anda silakan klik URL ini.</p><p><a href=\"http://" + config.SERVER_HOST + ":" + config.SERVER_PORT + "/aktifkan-akun?token=" + activatehas.String +
			"\">Aktifkan akun</a>.</p><br><p>Jika anda tidak mendaftar pada website ini, silakan abaikan email ini.</p><div dir=\"ltr\">" +
			"<div dir=\"ltr\"><div><span style=\"font-size:12.8px\">Hormat Kami,</span></div><div>Telp./SMS/Whatsapp +6281311394295</div><div></div><div>" +
			"Email: pmb@stkipnu.ac.id<br><img src=\"https://pmb.stkipnu.ac.id/assets/images/logo-stkipnu.png\" height=\"45\" class=\"CToWUd\"><br></div></div></div>"

		helper.KirimAktifasi("Aktifasi Akun", htmlaktivasi, data.Email)
	}

	userModel := models.User{
		Username:       data.Username,
		Email:          data.Email,
		PasswordHash:   passwordHash,
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
