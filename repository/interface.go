package repository

import (
	"context"
	"github.com/server-catalog/models"
)

type CatalogRepository interface {
	Upload(ctx context.Context, servers []models.ServerCatalog) error
	GetLocations(ctx context.Context) ([]string, error)
	GetHDDTypes(ctx context.Context) ([]string, error)
}
