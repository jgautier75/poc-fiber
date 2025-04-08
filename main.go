package main

import (
	"fmt"
	"poc-fiber/logger"
	"poc-fiber/migrate"

	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic(fmt.Errorf("config file not found : [%w]", err))
		} else {
			// Config file was found but another error was produced
		}
	}

	logger := logger.ConfigureLogger(viper.GetString("app.logFile"), true, true)
	migrate.PerformMigration(logger, viper.GetString("app.pgAdminUrl"))
}
