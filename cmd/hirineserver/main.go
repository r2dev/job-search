package main

import (
	"hirine/app"
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config.dev")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %s \n", err)
	}
	app := app.CreateServer()
	app.Start()
}
