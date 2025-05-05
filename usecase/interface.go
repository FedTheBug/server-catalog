package usecase

import (
	"context"
	"github.com/server-catalog/models"
)

type CatalogInterface interface {
	UploadCatalog(ctx context.Context, ctr []models.ServerCatalog) error
}
