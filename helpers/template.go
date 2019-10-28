package helpers

import (
	"net/http"
	"sync"
	"text/template"

	"github.com/tj/go/http/response"
)

func HandleTemplate(files ...string) http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles(files...)
		})
		if err != nil {
			response.InternalServerError(w, "template error")
			return
		}
		tpl.Execute(w, nil)
	}
}
