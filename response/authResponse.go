package response

type (
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
