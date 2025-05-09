package v1

import (
	"tablelink/src/v1/handler"

	"github.com/go-chi/chi/v5"
)

func Router(r chi.Router, deps *Dependency) {
	// Check Health
	r.Get("/health", handler.Health())
}
