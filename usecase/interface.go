package usecase

import (
	"context"
	"github.com/server-catalog/internal/dto"
)

type CatalogUseCase interface {
	UploadCatalog(ctx context.Context, ctr *dto.UploadCatalogCtr) error
	GetLocations(ctx context.Context) ([]string, error)
	GetHDDTypes(ctx context.Context) ([]string, error)
	GetListOfServers(ctx context.Context, ctr *dto.ListServersCtr) ([]dto.ListServerResp, error)
}
