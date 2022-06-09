package handlers

import (
	"auth/V1/Auth/dto"
	"auth/V1/Auth/services"
	"auth/config"
	"auth/helper"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
	authService services.AuthService
}

func AuthHandler(authService services.AuthService) *authHandler {
	return &authHandler{authService}
}

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

type IError struct {
	Field string
	Tag   string
	Value string
}

func (h *authHandler) Register(c echo.Context) (err error) {

	auth := new(dto.AuthRegReq)

	if err = c.Bind(auth); err != nil {
		return
	}
	id := id.New()
	uni = ut.New(id, id)
	validate = validator.New()
	trans, _ := uni.GetTranslator("id")
	id_translations.RegisterDefaultTranslations(validate, trans)

	if auth.Password != auth.PasswordConfirm {
		responsError := dto.AuthResponse{
			Password: "Password tidak sama",
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"email":    responsError.Email,
			"username": responsError.Username,
			"password": responsError.Password,
			"status":   false,
		})
	}

	if err = validate.Struct(auth); err != nil {

		errs := err.(validator.ValidationErrors)
		errors := errs.Translate(trans)

		responsError := dto.AuthResponse{
			Email:    errors["AuthReq.Email"],
			Username: errors["AuthReq.Username"],
			Password: errors["AuthReq.Password"],
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"email":    responsError.Email,
			"username": responsError.Username,
			"password": responsError.Password,
			"error":    "Gagal mendaftar",
			"status":   false,
		})
	}

	IsRegistered, _ := h.authService.IsRegistered(*auth)
	var mailError = ""
	var usernameError = ""

	if IsRegistered.ID != 0 {
		if IsRegistered.Email == auth.Email {
			mailError = "Email sudah terdaftar !"
		}

		if IsRegistered.Username == auth.Username {
			usernameError = "Username sudah terdaftar !"
		}

		responsError := dto.AuthResponse{
			Email:    mailError,
			Username: usernameError,
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"email":    responsError.Email,
			"username": responsError.Username,
			"password": responsError.Password,
			"pesan":    "Gagal mendaftar",
			"status":   false,
		})
	}

	savedUser, _ := h.authService.Save(*auth)

	if savedUser.ID == 0 {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"pesan":  "Internal server error",
			"status": false,
		})
	}

	responsError := dto.AuthResponse{
		Email:    savedUser.Email,
		Username: savedUser.Username,
		Password: "-Rahasia-",
	}

	return c.JSON(http.StatusOK, echo.Map{
		"email":    responsError.Email,
		"username": responsError.Username,
		"password": responsError.Password,
		"pesan":    "Berhasil mendaftar",
		"status":   true,
	})
}

func (h *authHandler) Login(c echo.Context) (err error) {

	auth := new(dto.AuthLoginReq)

	if err = c.Bind(auth); err != nil {
		return
	}

	IsRegistered, _ := h.authService.IsRegisteredForLogin(*auth)

	if !helper.CheckPasswordHash(auth.Password, IsRegistered.PasswordHash) {
		b := &dto.AuthLoginReq{
			IPAddress: c.RealIP(),
			UserID:    IsRegistered.ID,
			Success:   0,
		}

		h.authService.SaveAuthLogin(*b)
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"pesan":  "Username atau password salah",
			"status": false,
		})
	}

	b := &dto.AuthLoginReq{
		IPAddress: c.RealIP(),
		UserID:    IsRegistered.ID,
		Success:   1,
	}

	h.authService.SaveAuthLogin(*b)

	claims := &helper.JwtCustomClaims{
		Username:   IsRegistered.Username,
		UserID:     IsRegistered.ID,
		Group:      "admin",
		Permission: "Full Controll",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * config.JWT_EXPIRE).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.JWT_SECREET_KEY))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":  t,
		"pesan":  "Berhasil masuk",
		"status": true,
	})
}

func (h *authHandler) AddGroup(c echo.Context) (err error) {

	auth := new(dto.AuthGroupReq)

	if err = c.Bind(auth); err != nil {
		return
	}

	id := id.New()
	uni = ut.New(id, id)
	validate = validator.New()
	trans, _ := uni.GetTranslator("id")
	id_translations.RegisterDefaultTranslations(validate, trans)

	if err = validate.Struct(auth); err != nil {

		errs := err.(validator.ValidationErrors)
		errors := errs.Translate(trans)

		responsError := dto.GroupResponse{
			Name: errors["AuthGroup.Name"],
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"name":   responsError.Name,
			"pesan":  "Gagal memasukan data group",
			"status": false,
		})
	}

	h.authService.AddGroup(*auth)

	return c.JSON(http.StatusOK, echo.Map{
		// "token":  t,
		"pesan":  "Berhasil memasukan data group",
		"status": true,
	})
}

func (h *authHandler) AddPermission(c echo.Context) (err error) {

	auth := new(dto.AuthPermissionReq)

	if err = c.Bind(auth); err != nil {
		return
	}

	id := id.New()
	uni = ut.New(id, id)
	validate = validator.New()
	trans, _ := uni.GetTranslator("id")
	id_translations.RegisterDefaultTranslations(validate, trans)

	if err = validate.Struct(auth); err != nil {

		errs := err.(validator.ValidationErrors)
		errors := errs.Translate(trans)

		responsError := dto.GroupResponse{
			Name: errors["AuthGroup.Name"],
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"name":   responsError.Name,
			"pesan":  "Gagal memasukan data permission",
			"status": false,
		})
	}

	h.authService.AddPermission(*auth)

	return c.JSON(http.StatusOK, echo.Map{
		// "token":  t,
		"pesan":  "Berhasil memasukan data permission",
		"status": true,
	})
}

func (h *authHandler) GetUsers(c echo.Context) (err error) {

	// fmt.Println("Ofset : ", c.QueryParam("start"))

	// fmt.Println(c.FormValue("order[0][column]"))
	// fmt.Println(c.FormValue("order[0][dir]"))

	auth := new(dto.DatatablesReq)

	if err = c.Bind(auth); err != nil {
		return
	}

	fmt.Println(auth)

	// users, _ := h.authService.FindAll()
	q, users, _ := h.authService.Datatables(*auth)

	var res []*dto.UserResponse
	var i uint
	for _, user := range users {
		i++
		var el dto.UserResponse
		el.ID = i
		el.Username = user.Username
		el.Email = user.Email
		el.Active = user.Active
		var al []string
		for _, name := range user.AuthGroupUser {
			al = append(al, name.AuthGroup.Name)
		}
		el.Group = al
		res = append(res, &el)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":            res,
		"draw":            auth.Draw,
		"recordsFiltered": q,
		"recordsTotal":    11,
		"username":        "Berhasil mendapatkan data user",
		"email":           "hello",
		"active":          true,
	})
}

// func validMailAddress(address string) (string, bool) {
// 	addr, err := mail.ParseAddress(address)
// 	if err != nil {
// 		return "", false
// 	}
// 	return addr.Address, true
// }
