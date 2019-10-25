package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type LoginWithPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginWithPassword(c echo.Context) error {
	request := new(LoginWithPasswordRequest)
	if err := c.Bind(request); err != nil {
		// @todo
		return echo.ErrBadGateway
	}
	username := request.Username
	password := request.Password

	// Throws unauthorized error
	if username != "jon" || password != "shhh!" {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Jon Snow"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
