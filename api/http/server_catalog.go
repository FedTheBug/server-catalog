package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/server-catalog/internal/dto"
	"github.com/server-catalog/internal/utils"
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

		r.Get("/servers/hdd-types", handler.getHddTypes)
		r.Get("/servers/locations", handler.getLocations)

		r.Get("/servers/list", handler.getServers)
	})
}

func (s *SCHandler) getServers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page := utils.NewPage(r)

	data, err := s.scUseCase.GetListOfServers(ctx, &dto.ListServersCtr{Page: page})
	if err != nil {
		_ = (&utils.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "unable to fetch servers",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	_ = (&utils.Response{
		Status:     http.StatusOK,
		Pagination: page,
		Data:       data,
	}).Render(w)

	return
}

func (s *SCHandler) getHddTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := s.scUseCase.GetHDDTypes(ctx)
	if err != nil {
		_ = (&utils.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "unable to fetch hdd types",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	_ = (&utils.Response{
		Status: http.StatusOK,
		Data:   data,
	}).Render(w)

	return
}

func (s *SCHandler) getLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := s.scUseCase.GetLocations(ctx)
	if err != nil {
		_ = (&utils.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "unable to fetch locations",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	_ = (&utils.Response{
		Status: http.StatusOK,
		Data:   data,
	}).Render(w)

	return
}

func (s *SCHandler) uploadCatalog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		_ = (&utils.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "unable to parse form",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		_ = (&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "file is required",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	if err := s.scUseCase.UploadCatalog(ctx, &dto.UploadCatalogCtr{File: file}); err != nil {
		_ = (&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "failed to upload file",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	defer file.Close()

	_ = (&utils.Response{
		Status:  http.StatusCreated,
		Message: "Catalog uploaded",
		Error:   err,
	}).Render(w)

	return
}
