package main

import (
	"fmt"
	"hirine/auth"
	"hirine/models"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// func accessible(c echo.Context) error {
// 	return c.String(http.StatusOK, "Accessible")
// }

// func restricted(c echo.Context) error {
// 	user := c.Get("user").(*jwt.Token)
// 	claims := user.Claims.(jwt.MapClaims)
// 	name := claims["name"].(string)
// 	return c.String(http.StatusOK, "Welcome "+name+"!")
// }

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
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	models.InitMongo(viper.GetString("mongo_url"))
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/auth/register", auth.RegisterWithPassword)
	r.Post("/auth/login-username", auth.LoginWithPassword)
	// e.POST("/auth/verify-phone", auth.VerifyPhone)
	// e.POST("/auth/change-password", auth.ChangePassword)

	r.Post("/auth/register-phone", auth.RegisterWithPhoneNumber)

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
