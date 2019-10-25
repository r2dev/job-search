package job

import (
	"net/http"

	"github.com/labstack/echo"
)

func Create(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}
