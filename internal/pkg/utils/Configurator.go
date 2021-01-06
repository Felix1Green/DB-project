package utils

import (
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/spf13/viper"
	"log"
)

type ServiceConfig struct {
	Domain           string
	Port             int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseDomain string
	DatabasePort   int
}


func Run(configPath string) (*ServiceConfig, error) {
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Unable to read config file: %s", err)
		return nil, models.IncorrectPath
	}

	config := new(ServiceConfig)
	config.Domain = viper.GetString("Domain")
	config.Port = viper.GetInt("Port")
	config.DatabaseName = viper.GetString("Database.Name")
	config.DatabaseUser = viper.GetString("Database.User")
	config.DatabaseDomain = viper.GetString("Database.Domain")
	config.DatabasePort = viper.GetInt("Database.Port")
	config.DatabasePassword = viper.GetString("Database.Password")
	return config, nil
}