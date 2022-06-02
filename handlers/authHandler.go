package handlers

import (
	"auth/requests"
	"auth/response"
	"auth/services"
	"net/http"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
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

	if err = validate.Struct(auth); err != nil {
		errs := err.(validator.ValidationErrors)
		errors := errs.Translate(trans)

		responsError := response.AuthResponse{
			Email:    errors["AuthReq.Email"],
			Username: errors["AuthReq.Username"],
			Password: errors["AuthReq.Password"],
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"err": responsError,
		})
	}

	s, er := h.authService.Save(*auth)

	return c.JSON(http.StatusOK, echo.Map{
		"auth": s,
		"err":  er,
	})
}

func (h *authHandler) Login(c echo.Context) (err error) {

	auth := new(requests.AuthReq)

	if err = c.Bind(auth); err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"auth": auth,
	})
}
