package dto

type (
	AuthRegReq struct {
		Email           string `json:"email" validate:"required,email"`
		Username        string `json:"username" validate:"required,alphanum,lowercase,min=3,max=30"`
		Password        string `json:"password" validate:"required,min=8"`
		PasswordConfirm string `json:"password_confirm" validate:"required,min=8"`
	}

	AuthLoginReq struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		IPAddress string
		UserID    uint
		Success   int
	}

	AuthGroupReq struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	AuthPermissionReq struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}
)

type (
	UserResponse struct {
		ID       uint        `json:"id"`
		Email    string      `json:"email"`
		Username string      `json:"username"`
		Active   bool        `json:"active"`
		Group    interface{} `json:"role"`
	}

	AuthResponse struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	GroupResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)
