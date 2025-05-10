package migration

import (
	"fmt"
	"github.com/server-catalog/internal/conn"
	"github.com/spf13/cobra"
)

var (

	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:   "migration",
		Short: "Run database migrations",
		Long:  `Migration is a tool to generate and modify database tables`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := conn.ConnectDB(); err != nil {
				return fmt.Errorf("cant't connect database: %v", err)
			}
			return nil
		},
	}
)
