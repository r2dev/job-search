package company

import (
	"net/http"

	"github.com/tj/go/http/response"
)

func GetCompany(w http.ResponseWriter, r *http.Request) {
	response.OK(w, 1)
}
