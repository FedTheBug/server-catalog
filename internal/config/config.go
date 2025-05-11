package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var mu sync.Mutex

// LoadConfig loads configurations from config.yml file
func LoadConfig() error {
	viper.SetConfigFile("config.example.yml") // local

	if err := viper.ReadInConfig(); err != nil {
		fmt.Errorf("error reading config file: %s", err)
	}

	LoadApp()
	LoadDB()

	return nil
}
