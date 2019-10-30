package app

import (
	"hirine/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type App struct {
	DB *models.DB
	R  *chi.Mux
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
	csrfMiddleware := csrf.Protect([]byte(viper.GetString("session_secret")), csrf.Secure(false))

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "templates/static")
	fileServer(r, "/static", http.Dir(filesDir))
	r.Group(func(r chi.Router) {
		r.Use(csrfMiddleware)
		r.Get("/", app.IndexPage)
		r.Get("/register", app.RegisterPage)
		r.Post("/register", app.RegisterHandler)
		r.Get("/login", app.LoginPage)
		r.Post("/login", app.LoginHandler)
	})

	r.Post("/auth/register", app.RegisterWithPassword)
	r.Post("/auth/login-username", app.LoginWithPassword)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/company", app.GetCompany)
		r.Post("/company", app.CreateCompany)
		r.Put("/company/{id}", app.UpdateCompany)
		r.Delete("/company/{id}", app.DeleteCompany)

		r.Post("/job", app.CreateJob)
		r.Put("/job/{id}", app.UpdateJob)
		r.Delete("/job/{id}", app.DeleteJob)

		r.Post("/application/apply", app.ApplyJob)
	})

	// log.Info("starting server on 1323")
	// log.Fatal(http.ListenAndServe(":1323", r))
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
