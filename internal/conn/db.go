package conn

import (
	"fmt"
	"github.com/server-catalog/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// DB holds the database instance
var (
	db *gorm.DB
)

// Ping tests if db connection is alive
func Ping() error {
	d, _ := db.DB() // returns *sql.DB
	return d.Ping()
}

var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             500 * time.Millisecond, // Slow SQL threshold
		LogLevel:                  logger.Info,            // Log level
		IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
		Colorful:                  true,                   // Disable color
	},
)

// Connect sets the db client of database using configuration cfg
func Connect(cfg *config.Database) error {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	gormDB, err := gorm.Open(mysql.Open(uri), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		return err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	if cfg.MaxIdleConnection != 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConnection)
	}
	if cfg.MaxActiveConnection != 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxActiveConnection)
	}

	db = gormDB

	return nil
}

// DefaultDB returns default db
func DefaultDB() *gorm.DB {
	if Ping() != nil { // if connection lost, reconnect it
		if err := ConnectDB(); err != nil {
			log.Println(err)
		}
	}
	if config.App().Env == config.EnvDevelopment ||
		config.App().Env == config.EnvStaging {
		db.Debug()
	}
	return db
}

// ConnectDB sets the db client of database using default configuration file
func ConnectDB() error {
	cfg := config.DB()
	return Connect(cfg)
}
