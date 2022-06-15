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

	DatatablesReq struct {
		Length       int    `form:"length"`
		Draw         int    `form:"draw"`
		Start        int    `form:"start"`
		Order0Column int    `form:"order[0][column]"`
		Order0Dir    string `form:"order[0][dir]"`
		SearchValue  string `form:"search[value]"`
	}
)

type (
	UserResponse struct {
		ID       uint        `json:"no"`
		Email    string      `json:"email"`
		Username string      `json:"username"`
		Active   bool        `json:"status"`
		Group    interface{} `json:"role"`
		UserID   uint        `json:"user_id"`
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
