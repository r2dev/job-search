package job

import (
	"net/http"

	"github.com/labstack/echo"
)

func Delete(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}
