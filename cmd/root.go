package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	RootCmd = &cobra.Command{
		Use:   "Golang-restful Server Catalog",
		Short: "A http service",
		Long:  `An HTTP JSON API backend service`,
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
