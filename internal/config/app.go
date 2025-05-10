package config

import "github.com/spf13/viper"

// represents environment level
const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

// Application represents application config
type Application struct {
	Base            string
	Port            int
	Env             string
	PaginationLimit int
	SecretKey       string
}

var app Application

// App contains app configurations
func App() *Application {
	return &app
}

func LoadApp() {
	mu.Lock()
	defer mu.Unlock()

	app = Application{
		Base:            viper.GetString("app.host"),
		Port:            viper.GetInt("app.port"),
		Env:             viper.GetString("app.env"),
		PaginationLimit: viper.GetInt("app.pagination_limit"),
		SecretKey:       viper.GetString("app.secret_key"),
	}
}
