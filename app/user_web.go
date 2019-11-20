package app

import (
	"net/http"
	"text/template"

	"github.com/gorilla/csrf"
)

func (app *App) RegisterUserGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		}
		t := template.Must(
			template.ParseFiles("./templates/layout/base.html", "./templates/register.html"))

		t.Execute(w, map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
			"messages":       messages,
		})
	}
}

func (app *App) LoginUserGet(w http.ResponseWriter, r *http.Request) {
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
	}
	var indexTemp = template.Must(
		template.ParseFiles("./templates/layout/base.html", "./templates/login.html"))
	indexTemp.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"messages":       messages,
	})
}
