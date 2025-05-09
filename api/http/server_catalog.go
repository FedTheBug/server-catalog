package http

import (
	"errors"
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

	var storageMin, storageMax *int
	if minStr := r.URL.Query().Get("min_storage"); minStr != "" {
		if min, err := utils.ParseStorageToGB(minStr); err == nil {
			storageMin = &min
		}
	}
	if maxStr := r.URL.Query().Get("max_storage"); maxStr != "" {
		if max, err := utils.ParseStorageToGB(maxStr); err == nil {
			storageMax = &max
		}
	}

	var ramValues []int
	if ramStr := r.URL.Query().Get("ram"); ramStr != "" {
		ramValues = utils.ParseRAMValues(ramStr)
	}

	var hddTypeID *int
	if hdd := r.URL.Query().Get("hdd_type"); hdd != "" {
		if id, err := utils.GetHDDTypeID(hdd); err == nil {
			hddTypeID = &id
		}
	}

	var location *string
	if loc := r.URL.Query().Get("location"); loc != "" {
		location = &loc
	}

	data, err := s.scUseCase.GetListOfServers(ctx, &dto.ListServersCtr{
		StorageMin: storageMin,
		StorageMax: storageMax,
		RAM:        ramValues,
		HDD:        hddTypeID,
		Location:   location,
		Page:       page,
	})

	if err != nil {
		if errors.Is(err, utils.ErrServerNotFound) {
			_ = (&utils.Response{
				Status:  http.StatusNotFound,
				Message: "no server found with these configs",
				Error:   err.Error(),
			}).Render(w)
			return
		}
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
