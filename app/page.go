package app

import (
	"net/http"
	"text/template"

	"github.com/gorilla/csrf"
)

func (app *App) IndexPage(w http.ResponseWriter, r *http.Request) {
	session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
	_, ok := session.Values["n_0"]
	login := false
	if ok {
		login = true
	}
	indexTemp := template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/index.html"))
	indexTemp.Execute(w, map[string]interface{}{
		"login":          login,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (app *App) RegisterPage(w http.ResponseWriter, r *http.Request) {
	session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
	_, ok := session.Values["n_0"]
	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t := template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/register.html"))

	t.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (app *App) LoginPage(w http.ResponseWriter, r *http.Request) {
	session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
	_, ok := session.Values["n_0"]
	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var indexTemp = template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/login.html"))
	indexTemp.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
