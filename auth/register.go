package auth

import (
	"net/http"

	"github.com/labstack/echo"
)

func Register(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"token": "efwfw",
	})
}
