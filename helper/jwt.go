package helper

import "github.com/golang-jwt/jwt"

type jwtCustomClaims struct {
	Username   string `json:"username"`
	UserID     uint   `json:"user_id"`
	Group      string `json:"group"`
	Permission string `json:"permission"`
	jwt.StandardClaims
}

// func login(c echo.Context) error {
// 	username := c.FormValue("username")
// 	password := c.FormValue("password")

// 	// Throws unauthorized error
// 	if username != "jon" || password != "shhh!" {
// 		return echo.ErrUnauthorized
// 	}

// 	// Set custom claims
// 	claims := &jwtCustomClaims{
// 		"Jon Snow",
// 		"jonsnow",
// 		true,
// 		jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
// 		},
// 	}

// 	// Create token with claims
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	// Generate encoded token and send it as response.
// 	t, err := token.SignedString([]byte("secret"))
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusOK, echo.Map{
// 		"token": t,
// 	})
// }

// func accessible(c echo.Context) error {
// 	return c.String(http.StatusOK, "Accessible")
// }

// func restricted(c echo.Context) error {
// 	user := c.Get("user").(*jwt.Token)
// 	claims := user.Claims.(*jwtCustomClaims)
// 	name := claims.Name
// 	admin := claims.Admin
// 	return c.String(http.StatusOK, "Welcome "+name+"!"+"hh"+strconv.FormatBool(admin))
// }
