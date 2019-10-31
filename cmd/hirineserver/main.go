package main

import (
	"hirine/app"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	viper.SetConfigName("config.dev")
	viper.AddConfigPath("./config")
}

func main() {

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %s \n", err)
	}
	app := app.CreateServer()
	app.Start()
}
