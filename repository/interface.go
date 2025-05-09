package repository

import (
	"context"
	"github.com/server-catalog/internal/dto"
	"github.com/server-catalog/models"
)

type CatalogRepository interface {
	Upload(ctx context.Context, servers []models.ServerCatalog) error
	GetLocations(ctx context.Context) ([]string, error)
	GetHDDTypes(ctx context.Context) ([]string, error)
	GetServers(ctx context.Context, ctr *dto.ListServersCtr) ([]models.ServerCatalog, error)
}
