package handler

import (
	"net/http"
	"text/template"

	"github.com/gorilla/csrf"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	indexTemp := template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/index.html"))
	indexTemp.Execute(w, nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/register.html"))

	t.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	var indexTemp = template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/login.html"))
	indexTemp.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
