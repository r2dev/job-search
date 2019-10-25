package job

import (
	"net/http"

	"github.com/labstack/echo"
)

func Update(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}
