package requests

type (
	AuthReq struct {
		Email           string `json:"email" validate:"required,email"`
		Username        string `json:"username" validate:"required,alphanum,lowercase,min=3,max=30"`
		Password        string `json:"password" validate:"required,min=8"`
		PasswordConfirm string `json:"password_confirm" validate:"required,min=8"`
	}

	AuthLogin struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)
