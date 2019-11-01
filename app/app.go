package app

import (
	"hirine/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type App struct {
	DB *models.DB
	R  *chi.Mux
	S  *sessions.CookieStore
	E  *casbin.Enforcer
}

// CreateServer create a server instance
func CreateServer() *App {
	app := &App{}

	db, err := models.InitMongo(viper.GetString("mongo_url"))
	if err != nil {
		log.WithError(err).Fatal("config failed")
	}
	app.DB = db
	tokenAuth := jwtauth.New("HS256", []byte(viper.GetString("jwt_secret")), nil)
	csrfMiddleware := csrf.Protect([]byte(viper.GetString("csrf_secret")), csrf.Secure(false))
	store := sessions.NewCookieStore([]byte(viper.GetString("session_secret")))
	app.S = store
	e, _ := casbin.NewEnforcer(viper.GetString("casbin_model"), viper.GetString("casbin_policy"))
	app.E = e
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "templates/static")
	fileServer(r, "/static", http.Dir(filesDir))
	r.Group(func(r chi.Router) {
		r.Use(csrfMiddleware)
		r.Get("/", app.IndexGet())
		r.Get("/register", app.RegisterUserGet)
		r.Post("/register", app.RegisterUserPost)
		r.Get("/login", app.LoginUserGet)
		r.Post("/login", app.LoginUserPost)
		r.Post("/logout", app.LogoutUserPost)

		r.Get("/dashboard", app.DashboardGet())
		r.Get("/dashboard/company", app.DashboardCompanyGet())
		r.Get("/dashboard/company/{companyID}/admin", app.CompanyAdminGet())
		r.Get("/dashboard/company-register", app.RegisterCompanyGet())
		r.Post("/dashboard/company-register", app.RegisterCompanyPost())
		// r.Get("/dashboard/company/{companyID}/job", app.CompanyJobGet())

		// r.Post("/company/{companyID}/job", app.CompanyJobPost)

	})

	r.Post("/auth/register", app.RegisterWithPassword)
	r.Post("/auth/login-username", app.LoginWithPassword)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/api/company", app.GetCompany)
		r.Post("/api/company", app.CreateCompany)
		r.Put("/api/company/{id}", app.UpdateCompany)
		r.Delete("/api/company/{id}", app.DeleteCompany)

		r.Post("/api/job", app.CreateJob)
		r.Put("/api/job/{id}", app.UpdateJob)
		r.Delete("/api/job/{id}", app.DeleteJob)

		r.Post("/api/application/apply", app.ApplyJob)
	})

	app.R = r
	return app
}

// Start starts application
func (app *App) Start() {
	log.Fatal(http.ListenAndServe(":1323", app.R))
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
