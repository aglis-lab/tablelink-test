package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"tablelink/src/app"
	"tablelink/src/middleware"
	v1 "tablelink/src/v1"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Init app context
	if err := app.Init(ctx); err != nil {
		log.Panic(err)
	}

	// Init Router
	address := fmt.Sprintf(":%d", app.Config().BindAddress)
	r := initRouter(ctx)

	log.Println("Listening to", address)
	http.ListenAndServe(address, r)
}

func initRouter(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.RequestIDContext(middleware.DefaultGenerator))
	r.Use(middleware.RequestAttributesContext)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	deps := v1.Dependencies(ctx)

	r.Route("/v1", func(r chi.Router) {
		v1.Router(r, deps)
	})

	return r
}
