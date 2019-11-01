package app

import (
	"fmt"
	"net/http"
	"sync"
	"text/template"

	"github.com/tj/go/http/response"

	"github.com/gorilla/csrf"
)

func (app *App) HandleIndex() http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles("./templates/layout/base.html", "./templates/index.html")
		})
		if err != nil {
			response.InternalServerError(w, err.Error())
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		_, ok := session.Values["n_0"]
		login := false
		if ok {
			login = true
		}
		tpl.Execute(w, map[string]interface{}{
			"login":          login,
			csrf.TemplateTag: csrf.TemplateField(r),
		})
	}

}

func (app *App) RegisterPage(w http.ResponseWriter, r *http.Request) {
	session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
	flash := session.Flashes()
	session.Save(r, w)
	var messages []string
	if _, ok := session.Values["n_0"]; ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if flash != nil {
		for _, f := range flash {
			fString, ok := f.(string)
			if ok {
				messages = append(messages, fString)
			}

		}
	} else {
		fmt.Println("no flash message")
	}
	t := template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/register.html"))

	t.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"messages":       messages,
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

func (app *App) RegisterCompanyGet() http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles("./templates/layout/base.html", "./templates/company_register.html")
		})
		if err != nil {
			response.InternalServerError(w, err.Error())
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		_, ok := session.Values["n_0"]
		login := false
		if ok {
			login = true
		}
		tpl.Execute(w, map[string]interface{}{
			"login":          login,
			csrf.TemplateTag: csrf.TemplateField(r),
		})
	}
}
