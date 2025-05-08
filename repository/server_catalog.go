package repository

import (
	"context"
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
