package migration

import (
	"errors"
	"fmt"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/server-catalog/internal/config"
	"github.com/spf13/cobra"
	"log"
)

var migrateCmd = &cobra.Command{
	Use:   "up",
	Short: "Up database migrations",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := config.LoadConfig(); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.DB()
		uri := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?parseTime=True", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
		migrationsPath := "file://db/migrations"

		m, err := migrate.New(migrationsPath, uri)
		if err != nil {
			log.Fatalf("Failed to create migrate instance: %v", err)
		}
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations applied successfully!")
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
