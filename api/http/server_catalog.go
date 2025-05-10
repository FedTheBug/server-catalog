package http

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/server-catalog/docs" // This will import the generated docs
	"github.com/server-catalog/internal/dto"
	"github.com/server-catalog/internal/utils"
	"github.com/server-catalog/middleware"
	"github.com/server-catalog/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type SCHandler struct {
	scUseCase usecase.CatalogUseCase
}

func New(router *chi.Mux, cuc usecase.CatalogUseCase) {
	handler := &SCHandler{
		scUseCase: cuc,
	}

	// Add CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "App-key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AppKeyResolver)
		r.Post("/upload", handler.uploadCatalog)

		r.Get("/servers/hdd-types", handler.getHddTypes)
		r.Get("/servers/locations", handler.getLocations)

		r.Get("/servers/list", handler.getServers)

	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

}

// @Summary      Get list of servers
// @Description  Retrieve a paginated list of servers with optional filtering by storage, RAM, HDD type, and location
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        per_page query int false "Number of items per page (default: 10)"
// @Param        page_no query int false "Page number (default: 1)"
// @Param        min_storage query string false "Minimum storage (e.g., 1TB)"
// @Param        max_storage query string false "Maximum storage (e.g., 100TB)"
// @Param        ram query string false "RAM values (e.g., 2GB,4GB)"
// @Param        hdd_type query string false "HDD type (e.g., SATA2, SAS, SSD)"
// @Param        location query string false "Server location (e.g., AmsterdamAMS-01)"
// @Security     AppKeyAuth
// @Success      200  {object}  utils.Response{data=[]dto.ListServerResp,pagination=utils.Page} "List of servers with pagination"
// @Failure      404  {object}  utils.Response{message=string,error=string} "No servers found with the specified filters"
// @Failure      422  {object}  utils.Response{message=string,error=string} "Unable to fetch servers"
// @Example      {data} [{"model":"HP DL120G7Intel G850","ram":"4GBDDR3","hdd":"4x1TBSATA2","location":"AmsterdamAMS-01","price":"€39.99"}]
// @Example      {pagination} {"per_page":10,"page_no":1,"total":486}
// @Router       /servers/list [get]
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

// @Summary      Get HDD types
// @Description  Retrieve a list of all available HDD types in the server catalog
// @Tags         servers
// @Accept       json
// @Produce      json
// @Security     AppKeyAuth
// @Success      200  {object}  utils.Response{data=[]string} "List of HDD types (e.g., SAS, SATA2, SSD)"
// @Failure      422  {object}  utils.Response{message=string,error=string} "Unable to fetch HDD types"
// @Example      {data} ["SAS", "SATA2", "SSD"]
// @Router       /servers/hdd-types [get]
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

// @Summary      Get server locations
// @Description  Retrieve a list of all available server locations in the catalog
// @Tags         servers
// @Accept       json
// @Produce      json
// @Security     AppKeyAuth
// @Success      200  {object}  utils.Response{data=[]string} "List of server locations"
// @Failure      422  {object}  utils.Response{message=string,error=string} "Unable to fetch locations"
// @Example      {data} ["AmsterdamAMS-01", "DallasDAL-10", "FrankfurtFRA-10", "Hong KongHKG-10", "San FranciscoSFO-12", "SingaporeSIN-11", "Washington D.C.WDC-01"]
// @Router       /servers/locations [get]
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

// @Summary      Upload server catalog
// @Description  Upload a server catalog file in XLSX format. The file must contain valid server catalog data.
// @Tags         servers
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Server catalog file (XLSX format)"
// @Security     AppKeyAuth
// @Success      201  {object}  utils.Response{message=string} "Catalog uploaded successfully"
// @Failure      400  {object}  utils.Response{message=string,error=string} "Invalid file format or upload failed"
// @Failure      422  {object}  utils.Response{message=string,error=string} "Unable to process the file"
// @Example      {file} "servers_filters_assignment.xlsx"
// @Router       /upload [post]
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
