package main

import (
	"hirine/app"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "timestamp",
			log.FieldKeyLevel: "severity",
			log.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	viper.SetConfigName("config.dev")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %s \n", err)
	}
	app := app.CreateServer()
	app.Start()
}
