package main

import (
	"hirine/handler"
	"hirine/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/gorilla/csrf"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
	viper.SetConfigName("config.dev")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Panicf("Fatal error config file: %s \n", err)
	}
	models.InitMongo(viper.GetString("mongo_url"))

	tokenAuth := jwtauth.New("HS256", []byte(viper.GetString("jwt_secret")), nil)
	csrfMiddleware := csrf.Protect([]byte(viper.GetString("session_secret")), csrf.Secure(false))

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "templates/static")
	FileServer(r, "/static", http.Dir(filesDir))
	r.Group(func(r chi.Router) {
		r.Use(csrfMiddleware)
		r.Get("/", handler.IndexPage)
		r.Get("/register", handler.RegisterPage)
		r.Post("/register", handler.RegisterHandler)
		r.Get("/login", handler.LoginPage)
		r.Post("/login", handler.LoginHandler)
	})

	r.Post("/auth/register", handler.RegisterWithPassword)
	r.Post("/auth/login-username", handler.LoginWithPassword)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/company", handler.GetCompany)
		r.Post("/company", handler.CreateCompany)
		r.Put("/company/{id}", handler.UpdateCompany)
		r.Delete("/company/{id}", handler.DeleteCompany)

		r.Post("/job", handler.CreateJob)
		r.Put("/job/{id}", handler.UpdateJob)
		r.Delete("/job/{id}", handler.DeleteJob)

		r.Post("/application/apply", handler.ApplyJob)
	})

	log.Info("starting server on 1323")
	log.Fatal(http.ListenAndServe(":1323", r))
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
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
