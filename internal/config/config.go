package config

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

var mu sync.Mutex

// LoadConfig loads configurations from config.yml file
func LoadConfig() error {
	viper.SetConfigFile("config.example.yml") // local

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	LoadApp()
	LoadDB()

	return nil
}
