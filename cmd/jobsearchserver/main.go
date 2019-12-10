package main

import (
	"github.com/r2dev/job-search/app"
	"log"

	"github.com/spf13/viper"
)

func init() {

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
