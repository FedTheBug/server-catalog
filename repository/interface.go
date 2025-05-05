package repository

import (
	"context"
	"github.com/server-catalog/models"
)

type CatalogRepository interface {
	Upload(ctx context.Context, servers []models.ServerCatalog) error
}
