package main

import (
	"fmt"
	"hirine/auth"
	"hirine/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {

	viper.SetConfigName("config.dev")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	models.InitMongo(viper.GetString("mongo_url"))

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Login route
	e.POST("/auth/login-username", auth.LoginWithPassword)
	// e.POST("/auth/verify-phone", auth.VerifyPhone)
	// e.POST("/auth/change-password", auth.ChangePassword)
	e.POST("/auth/register", auth.Register)

	// Unauthenticated route
	e.GET("/", accessible)

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

	e.Logger.Fatal(e.Start(":1323"))
}
