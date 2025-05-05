package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/server-catalog/usecase"
	"net/http"
)

type SCHandler struct {
	scUseCase usecase.CatalogUseCase
}

func New(router *chi.Mux, cuc usecase.CatalogUseCase) {
	handler := &SCHandler{
		scUseCase: cuc,
	}

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/upload", handler.uploadCatalog)
	})
}

func (s *SCHandler) uploadCatalog(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, world!"))
}
