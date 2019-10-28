package handler

import (
	"net/http"
	"text/template"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	var indexTemp = template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/index.html"))
	indexTemp.Execute(w, nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	var indexTemp = template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/register.html"))
	indexTemp.Execute(w, nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	var indexTemp = template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/login.html"))
	indexTemp.Execute(w, nil)
}
