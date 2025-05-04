package cmd

import (
	"context"
	"errors"
	"fmt"
	chi "github.com/go-chi/chi/v5"
	"github.com/server-catalog/internal/config"
	"github.com/server-catalog/internal/conn"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve starts the http server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := config.LoadConfig(); err != nil {
			log.Fatalln(err)
		}

		if err := conn.ConnectDB(); err != nil {
			log.Fatalln(err)
		}

	},
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {
	r := chi.NewRouter()

	hServer := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.App().Base, config.App().Port),
		Handler: r,
	}

	//catalog repo
	//TODO:: connect to catalog repo

	//TODO:: initiate usecase
	//TODO:: initiate handler

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("HTTP:: Listening on port ", config.App().Port)
		if err := hServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-stop

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	hServer.Shutdown(ctx)
}
