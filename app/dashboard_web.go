package app

import (
	"hirine/models"
	"net/http"
	"sync"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
	"github.com/tj/go/http/response"
)

func (app *App) DashboardGet() http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles(
				"./templates/layout/base-dashboard.html", "./templates/dashboard.html")
		})
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		login := false
		if _, ok := session.Values["n_0"]; ok {
			login = true
		}
		flash := session.Flashes()
		session.Save(r, w)
		var messages []string
		if flash != nil {
			for _, f := range flash {
				fString, ok := f.(string)
				if ok {
					messages = append(messages, fString)
				}

			}
		}
		tpl.Execute(w, map[string]interface{}{
			"login":          login,
			csrf.TemplateTag: csrf.TemplateField(r),
			"messages":       messages,
			"IsDashboard":    true,
		})
	}
}

func (app *App) DashboardCompanyGet() http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles(
				"./templates/layout/base-dashboard.html", "./templates/dashboard-company.html")
		})
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		login := false
		if _, ok := session.Values["n_0"]; ok {
			login = true
		}
		flash := session.Flashes()
		session.Save(r, w)
		var messages []string
		if flash != nil {
			for _, f := range flash {
				fString, ok := f.(string)
				if ok {
					messages = append(messages, fString)
				}

			}
		}
		tpl.Execute(w, map[string]interface{}{
			"login":          login,
			csrf.TemplateTag: csrf.TemplateField(r),
			"messages":       messages,
			"IsCompany":      true,
		})
	}
}

func (app *App) DashboardPostJobGet() http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles(
				"./templates/layout/base-dashboard.html", "./templates/post-job.html")
		})
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		login := false
		if _, ok := session.Values["n_0"]; ok {
			login = true
		}
		flash := session.Flashes()

		userID, ok := session.Values["n_0"].(string)
		if !ok {
			session.AddFlash("Please login first")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		session.Save(r, w)
		var messages []string
		if flash != nil {
			for _, f := range flash {
				fString, ok := f.(string)
				if ok {
					messages = append(messages, fString)
				}

			}
		}
		var companies []models.Company
		err := app.DB.GetCompaniesByAdminID(&companies, userID)
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		tpl.Execute(w, map[string]interface{}{
			"login":          login,
			csrf.TemplateTag: csrf.TemplateField(r),
			"messages":       messages,
			"companies":      companies,
		})
	}
}

func (app *App) DashboardJobListGet() http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles(
				"./templates/layout/base-dashboard.html", "./templates/dashboard-job-list.html",
			)
		})
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		userID, ok := session.Values["n_0"].(string)
		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		login := true
		flash := session.Flashes()
		session.Save(r, w)
		var messages []string
		if flash != nil {
			for _, f := range flash {
				fString, ok := f.(string)
				if ok {
					messages = append(messages, fString)
				}

			}
		}

		var jobs []models.Job
		err = app.DB.GetJobsByCreator(&jobs, userID)
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		tpl.Execute(w, map[string]interface{}{
			"login":          login,
			csrf.TemplateTag: csrf.TemplateField(r),
			"messages":       messages,
			"jobs":           jobs,
		})
	}
}

func (app *App) DashboardJobDetailGet() http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles(
				"./templates/layout/base-dashboard.html", "./templates/dashboard-job-detail.html")
		})
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		userID, ok := session.Values["n_0"].(string)
		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			app.L.WithField("user", userID)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		jobID := chi.URLParam(r, "jobID")
		var job models.Job
		err := app.DB.GetJobByID(&job, jobID)

		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}

		flash := session.Flashes()
		session.Save(r, w)
		var messages []string
		if flash != nil {
			for _, f := range flash {
				fString, ok := f.(string)
				if ok {
					messages = append(messages, fString)
				}
			}
		}
		tpl.Execute(w, map[string]interface{}{
			"login":          true,
			csrf.TemplateTag: csrf.TemplateField(r),
			"messages":       messages,
			"title":          job.Title,
			"id":             job.JobID.Hex(),
		})
	}
}

func (app *App) DashboardApplicationDetailGet() http.HandlerFunc {
	var (
		// init sync.Once
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		// init.Do(func() {
		tpl, err = template.ParseFiles(
			"./templates/layout/base-dashboard.html", "./templates/dashboard-application-detail.html")
		// })
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		_, ok := session.Values["n_0"].(string)

		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		login := true
		flash := session.Flashes()
		session.Save(r, w)
		var messages []string
		if flash != nil {
			for _, f := range flash {
				fString, ok := f.(string)
				if ok {
					messages = append(messages, fString)
				}
			}
		}
		applicationID := chi.URLParam(r, "applicationID")
		var application models.Application
		err := app.DB.GetApplicationByApplicationID(&application, applicationID)
		if err != nil {
			app.L.Errorln(err)
			response.InternalServerError(w, err.Error())
			return
		}
		var events []models.Event
		err = app.DB.GetEventsByApplicationID(&events, application.ApplicationID)
		tpl.Execute(w, map[string]interface{}{
			"login":          login,
			csrf.TemplateTag: csrf.TemplateField(r),
			"messages":       messages,
			"application":    application,
			"events":         events,
		})

	}
}

func (app *App) DashboardApplicationListGet() http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles(
				"./templates/layout/base-dashboard.html", "./templates/dashboard-job-application-list.html")
		})
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		_, ok := session.Values["n_0"].(string)

		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		login := true
		flash := session.Flashes()
		session.Save(r, w)
		var messages []string
		if flash != nil {
			for _, f := range flash {
				fString, ok := f.(string)
				if ok {
					messages = append(messages, fString)
				}
			}
		}
		jobID := chi.URLParam(r, "jobID")
		var applications []models.Application
		err := app.DB.GetApplicationsByJob(&applications, jobID)
		if err != nil {
			app.L.Errorln(err)
			response.InternalServerError(w, err.Error())
			return
		}
		tpl.Execute(w, map[string]interface{}{
			"login":          login,
			csrf.TemplateTag: csrf.TemplateField(r),
			"messages":       messages,
			"applications":   applications,
		})
	}
}
