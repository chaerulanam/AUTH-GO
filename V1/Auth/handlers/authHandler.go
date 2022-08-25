package handlers

import (
	"auth/V1/Auth/dto"
	"auth/V1/Auth/services"
	"auth/config"
	"auth/helper"
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

func (h *authHandler) Register(c echo.Context) (err error) {

	a := c.Request().Header
	api := a.Get("api-key")

	Api := h.authService.GetApi(api)
	if Api.ID != 0 {
		if Api.Expire.Unix() < time.Now().Unix() {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"meta": echo.Map{
					"code":    http.StatusBadRequest,
					"status":  "failed",
					"message": "API-KEY kadaluarsa",
				},
				"data": nil,
			})
		}
	} else {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"meta": echo.Map{
				"code":    http.StatusBadRequest,
				"status":  "failed",
				"message": "API-KEY tidak valid",
			},
			"data": nil,
		})
	}
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
		var msgErr = ""
		if len(responsError.Email) > 0 {
			msgErr = responsError.Email
		} else if len(responsError.Username) > 0 {
			msgErr = responsError.Username
		} else if len(responsError.Password) > 0 {
			msgErr = responsError.Password
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"meta": echo.Map{
				"code":    http.StatusBadRequest,
				"status":  "failed",
				"message": msgErr,
			},
			"data": nil,
		})
	}

	if err = validate.Struct(auth); err != nil {

		errs := err.(validator.ValidationErrors)
		errors := errs.Translate(trans)

		responsError := dto.AuthResponse{
			Email:    errors["AuthRegReq.Email"],
			Username: errors["AuthRegReq.Username"],
			Password: errors["AuthRegReq.Password"],
		}

		var msgErr = ""
		if len(responsError.Email) > 0 {
			msgErr = responsError.Email
		} else if len(responsError.Username) > 0 {
			msgErr = responsError.Username
		} else if len(responsError.Password) > 0 {
			msgErr = responsError.Password
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"meta": echo.Map{
				"code":    http.StatusBadRequest,
				"status":  "failed",
				"message": msgErr,
			},
			"data": nil,
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

		var msgErr = ""
		if len(responsError.Email) > 0 {
			msgErr = responsError.Email
		} else if len(responsError.Username) > 0 {
			msgErr = responsError.Username
		} else if len(responsError.Password) > 0 {
			msgErr = responsError.Password
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"meta": echo.Map{
				"code":    http.StatusBadRequest,
				"status":  "failed",
				"message": msgErr,
			},
			"data": nil,
		})
	}

	savedUser, _ := h.authService.Save(*auth)

	if len(config.DEFAULT_GROUP_NEW_USER) > 0 {
		h.authService.AddUserToGroup(savedUser.ID, config.DEFAULT_GROUP_NEW_USER)
	}

	if savedUser.ID == 0 {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"meta": echo.Map{
				"code":    http.StatusInternalServerError,
				"status":  "failed",
				"message": "Internal Server Error",
			},
			"data": nil,
		})
	}

	responsSuccess := dto.AuthResponse{
		Email:    savedUser.Email,
		Username: savedUser.Username,
		Password: "-Rahasia-",
	}

	return c.JSON(http.StatusOK, echo.Map{
		"meta": echo.Map{
			"code":    http.StatusOK,
			"status":  "success",
			"message": "Berhasil Mendaftar",
		},
		"data": echo.Map{
			"email":    responsSuccess.Email,
			"username": responsSuccess.Username,
			"password": responsSuccess.Password,
		},
	})
}

func (h *authHandler) Login(c echo.Context) (err error) {
	a := c.Request().Header
	api := a.Get("api-key")

	Api := h.authService.GetApi(api)
	if Api.ID != 0 {
		if Api.Expire.Unix() < time.Now().Unix() {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"meta": echo.Map{
					"code":    http.StatusBadRequest,
					"status":  "failed",
					"message": "API-KEY kadaluarsa",
				},
				"data": nil,
			})
		}
	} else {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"meta": echo.Map{
				"code":    http.StatusBadRequest,
				"status":  "failed",
				"message": "API-KEY tidak valid",
			},
			"data": nil,
		})
	}

	auth := new(dto.AuthLoginReq)

	if err = c.Bind(auth); err != nil {
		return
	}

	IsRegistered, _ := h.authService.IsRegisteredForLogin(*auth)

	if !helper.CheckPasswordHash(auth.Password, IsRegistered.PasswordHash) {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"meta": echo.Map{
				"code":    http.StatusBadRequest,
				"status":  "failed",
				"message": "Username atau password salah",
			},
			"data": nil,
		})
	}

	if !IsRegistered.Active {
		b := &dto.AuthLoginReq{
			IPAddress: c.RealIP(),
			UserID:    IsRegistered.ID,
			Success:   0,
		}
		h.authService.SaveAuthLogin(*b)

		return c.JSON(http.StatusUnauthorized, echo.Map{
			"meta": echo.Map{
				"code":    http.StatusBadRequest,
				"status":  "failed",
				"message": "Akun belum diaktifkan",
			},
			"data": nil,
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

	var al []string
	for _, name := range IsRegistered.AuthGroupUser {
		al = append(al, name.AuthGroup.Name)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"meta": echo.Map{
			"code":    http.StatusOK,
			"status":  "success",
			"message": "Berhasil masuk",
		},
		"data": echo.Map{
			"token":    t,
			"email":    IsRegistered.Email,
			"username": IsRegistered.Username,
			"user_id":  IsRegistered.ID,
			"role":     al,
		},
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

	auth := new(dto.DatatablesReq)

	if err = c.Bind(auth); err != nil {
		return
	}
	// users, _ := h.authService.FindAll()
	allCount, countFiltered, users, _ := h.authService.Datatables(*auth)

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
		el.UserID = user.ID
		res = append(res, &el)
	}

	if len(res) > 0 {
		return c.JSON(http.StatusOK, echo.Map{
			"data":            res,
			"draw":            auth.Draw,
			"recordsFiltered": countFiltered,
			"recordsTotal":    allCount,
			"pesan":           "Berhasil mendapatkan data",
			"status":          true,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data":            "",
		"draw":            auth.Draw,
		"recordsFiltered": 0,
		"recordsTotal":    0,
		"pesan":           "Gagal mendapatkan data",
		"status":          true,
	})
}
