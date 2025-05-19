package main

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"aiexec-plugin/internal/server"
	"aiexec-plugin/internal/types/app"
	"aiexec-plugin/internal/utils/log"
)

func main() {
	var config app.Config

	// load env
	godotenv.Load()

	err := envconfig.Process("", &config)
	if err != nil {
		log.Panic("Error processing environment variables: %s", err.Error())
	}

	config.SetDefault()

	if err := config.Validate(); err != nil {
		log.Panic("Invalid configuration: %s", err.Error())
	}

	(&server.App{}).Run(&config)
}
