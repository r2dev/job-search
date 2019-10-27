package main

import (
	"hirine/auth"
	"hirine/company"
	"hirine/job"
	"hirine/models"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/jwtauth"

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

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/auth/register", auth.RegisterWithPassword)
	r.Post("/auth/login-username", auth.LoginWithPassword)
	// e.POST("/auth/verify-phone", auth.VerifyPhone)
	// e.POST("/auth/change-password", auth.ChangePassword)

	// r.Post("/auth/register-phone", auth.RegisterWithPhoneNumber)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))

		r.Use(jwtauth.Authenticator)

		r.Get("/company", company.GetCompany)
		r.Post("/company", company.CreateCompany)
		r.Put("/company/{id}", company.UpdateCompany)
		r.Delete("/company/{id}", company.DeleteCompany)

		r.Post("/job", job.CreateJob)
		// r.Put("/job/{id}", company.UpdateJob)
		// r.
		// r.Post("/company/{id}/staff",)
	})

	// Unauthenticated route
	// r.Get("/", accessible)

	// r := e.Group("/employer")
	// r.Use(middleware.JWT([]byte("secret")))
	// r.POST("/company", company.Create)
	// r.PUT("/company/:id", company.Update)
	// r.POST("/company/:id/staff", company.AddStaff)

	// r.POST("/job", job.Create)
	// r.PUT("/job/:id", job.Update)
	// r.DELETE("/job/:id", job.SoftDelete)

	// r.POST("/job/:id/apply", application.Apply)
	// r.POST("/application/:id/send-interview", application.SendInterview)
	// r.POST("/application/:id/decline-interview", application.DeclineInterview)
	// r.POST("/application/:id/cancel-interview", application.CancelInterview)
	// r.POST("/application/:id/send-offer", application.SendOffer)
	// r.POST("/application/:id/cancel-offer", application.CancelOffer)

	// e.Logger.Fatal(e.Start(":1323"))
	log.Info("starting server on 1323")
	log.Fatal(http.ListenAndServe(":1323", r))
}
