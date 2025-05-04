package config

import "github.com/spf13/viper"

// Database holds the database configuration
type Database struct {
	Host                string
	Port                int
	Username            string
	Password            string
	Name                string
	MaxIdleConnection   int
	MaxActiveConnection int
}

var db Database

// DB returns the default database configuration
func DB() *Database {
	return &db
}

// LoadDB loads database configuration
func LoadDB() {
	mu.Lock()
	defer mu.Unlock()

	db = Database{
		Name:                viper.GetString("db.name"),
		Username:            viper.GetString("db.username"),
		Password:            viper.GetString("db.password"),
		Host:                viper.GetString("db.host"),
		Port:                viper.GetInt("db.port"),
		MaxIdleConnection:   viper.GetInt("db.max_idle_connections"),
		MaxActiveConnection: viper.GetInt("db.max_active_connections"),
	}
}
