package repository

import (
	"context"
	"fmt"
	"github.com/server-catalog/models"
	"gorm.io/gorm"
)

type ServerCatalog struct {
	db *gorm.DB
}

func NewServerCatalog(db *gorm.DB) CatalogRepository {
	return &ServerCatalog{db: db}
}

func (sc *ServerCatalog) Upload(ctx context.Context, servers []models.ServerCatalog) error {
	var tb models.ServerCatalog
	return sc.db.Table(tb.TableName()).Create(servers).Error
}

func (sc *ServerCatalog) GetLocations(ctx context.Context) ([]string, error) {
	var tb models.ServerCatalog
	locations := []string{}
	err := sc.db.Table(tb.TableName()).Select("DISTINCT location").Order("location").Pluck("location", &locations).Error
	if err != nil {
		return nil, fmt.Errorf("repository:usecase:: failed to fetch server locations %v", err.Error)
	}
	return locations, nil
}
