package handlers

import (
	"auth/config"
	"auth/helper"
	"auth/requests"
	"auth/response"
	"auth/services"
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

	auth := new(requests.AuthReq)

	if err = c.Bind(auth); err != nil {
		return
	}
	id := id.New()
	uni = ut.New(id, id)
	validate = validator.New()
	trans, _ := uni.GetTranslator("id")
	id_translations.RegisterDefaultTranslations(validate, trans)

	if auth.Password != auth.PasswordConfirm {
		responsError := response.AuthResponse{
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

		responsError := response.AuthResponse{
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

		responsError := response.AuthResponse{
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

	responsError := response.AuthResponse{
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

	auth := new(requests.AuthLogin)

	if err = c.Bind(auth); err != nil {
		return
	}

	IsRegistered, _ := h.authService.IsRegisteredForLogin(*auth)

	if !helper.CheckPasswordHash(auth.Password, IsRegistered.PasswordHash) {
		return echo.ErrUnauthorized
	}

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

func (h *authHandler) User(c echo.Context) (err error) {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*helper.JwtCustomClaims)
	username := claims.Username

	return c.JSON(http.StatusOK, echo.Map{
		"claim": username,
	})
}
