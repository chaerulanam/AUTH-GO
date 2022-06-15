package services

import (
	"auth/V1/Auth/dto"
	"auth/V1/Auth/repositories"
	"auth/config"
	"auth/helper"
	"auth/models"
	"database/sql"
)

type AuthService interface {
	FindAll() ([]models.User, error)
	Save(data dto.AuthRegReq) (models.User, error)
	IsRegistered(data dto.AuthRegReq) (models.User, error)
	IsRegisteredForLogin(data dto.AuthLoginReq) (models.User, error)
	SaveAuthLogin(data dto.AuthLoginReq) (models.AuthLogin, error)
	AddGroup(data dto.AuthGroupReq) (models.AuthGroup, error)
	AddPermission(data dto.AuthPermissionReq) (models.AuthPermission, error)
	Datatables(data dto.DatatablesReq) (int64, int64, []models.User, error)
	AddUserToGroup(user_id uint, name string) models.AuthGroupUser
}

func (s *authservice) AddUserToGroup(user_id uint, name string) models.AuthGroupUser {
	group, _ := s.userrepository.FindGroupId(name)

	b := models.AuthGroupUser{
		UserID:  user_id,
		GroupID: group.ID,
	}

	data, _ := s.userrepository.AddUserToGroup(b)
	return data
}

func (s *authservice) FindAll() ([]models.User, error) {
	user, err := s.userrepository.FindAll()
	return user, err
}

func (s *authservice) IsRegistered(data dto.AuthRegReq) (models.User, error) {

	userModel := models.User{
		Username: data.Username,
		Email:    data.Email,
	}

	user, err := s.userrepository.IsRegistered(userModel)
	return user, err
}

func (s *authservice) IsRegisteredForLogin(data dto.AuthLoginReq) (models.User, error) {
	userModel := models.User{
		Username: data.Username,
		Email:    data.Email,
	}
	user, err := s.userrepository.IsRegistered(userModel)
	return user, err
}

func (s *authservice) Save(data dto.AuthRegReq) (models.User, error) {

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

func (s *authservice) SaveAuthLogin(data dto.AuthLoginReq) (models.AuthLogin, error) {

	authloginModel := models.AuthLogin{
		IPAddress: data.IPAddress,
		UserID:    data.UserID,
		Success:   data.Success,
	}

	user, err := s.userrepository.SaveAuthLogin(authloginModel)
	return user, err
}

func (s *authservice) AddGroup(data dto.AuthGroupReq) (models.AuthGroup, error) {

	authgroupModel := models.AuthGroup{
		Name:        data.Name,
		Description: data.Description,
	}

	_data, err := s.userrepository.AddGroup(authgroupModel)
	return _data, err
}

func (s *authservice) AddPermission(data dto.AuthPermissionReq) (models.AuthPermission, error) {

	authpermissionModel := models.AuthPermission{
		Name:        data.Name,
		Description: data.Description,
	}

	_data, err := s.userrepository.AddPermission(authpermissionModel)
	return _data, err
}

func (s *authservice) Datatables(data dto.DatatablesReq) (int64, int64, []models.User, error) {
	// err := db.Model(&User).Count(&count).Error
	allcount, _ := s.userrepository.CountUsers()
	if len(data.SearchValue) > 0 {
		count, _ := s.userrepository.CountUserBySearch(data)
		user, err := s.userrepository.DatatablesSearch(data)
		return count, count, user, err
	} else {
		count := allcount
		user, err := s.userrepository.DatatablesFind(data)
		return allcount, count, user, err
	}

}

type authservice struct {
	userrepository repositories.AuthRepo
}

func AuthServ(userrepository repositories.AuthRepo) *authservice {

	return &authservice{userrepository}
}
