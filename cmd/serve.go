package cmd

import (
	"context"
	"errors"
	"fmt"
	chi "github.com/go-chi/chi/v5"
	cHttp "github.com/server-catalog/api/http"
	"github.com/server-catalog/internal/config"
	"github.com/server-catalog/internal/conn"
	"github.com/server-catalog/repository"
	"github.com/server-catalog/usecase"
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
		fmt.Println("Database Connected Successfully")

	},
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {
	r := chi.NewRouter()

	hServer := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.App().Base, config.App().Port),
		Handler: r,
	}

	// initialize repository
	catRepo := repository.NewServerCatalog(conn.DefaultDB())
	// initialize usecase
	catUseCase := usecase.New(catRepo)

	cHttp.New(r, catUseCase)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("╔════════════════════════════════════════════╗")
		log.Println("║             🚀 Server Starting             ║")
		log.Println("╠════════════════════════════════════════════╣")
		log.Println("HTTP:: Listening on port ", config.App().Port)
		log.Printf("║ 📚 API Docs: http://%s:%d/swagger/index.html\n", config.App().Base, config.App().Port)
		log.Println("╚════════════════════════════════════════════╝")
		if err := hServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-stop

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	hServer.Shutdown(ctx)
}
