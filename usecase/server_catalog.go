package usecase

import (
	"context"
	"github.com/server-catalog/models"
	"github.com/server-catalog/repository"
)

type ServerCatalog struct {
	SCRepo repository.CatalogRepository
}

func New(scr repository.CatalogRepository) CatalogUseCase {
	return &ServerCatalog{SCRepo: scr}
}

func (sc *ServerCatalog) UploadCatalog(ctx context.Context, ctr []models.ServerCatalog) error {
	return nil
}
