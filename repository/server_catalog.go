package repository

import (
	"context"
	"fmt"
	"github.com/server-catalog/internal/dto"
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
		return nil, fmt.Errorf("repository:server_catalog:: failed to fetch server locations %v", err)
	}
	return locations, nil
}

func (sc *ServerCatalog) GetHDDTypes(ctx context.Context) ([]string, error) {
	var hs models.HDDSpec
	types := []string{}
	err := sc.db.Table(hs.TableName()).Select("DISTINCT type").Order("type").Pluck("type", &types).Error
	if err != nil {
		return nil, fmt.Errorf("repository:server_catalog:: failed to fetch hdd specs %v", err)
	}
	return types, nil
}

func (sc *ServerCatalog) GetServers(ctx context.Context, ctr *dto.ListServersCtr) ([]models.ServerCatalog, error) {
	m := models.ServerCatalog{}
	res := []models.ServerCatalog{}
	qry := sc.db.Table(m.TableName())

	var count int64
	if err := qry.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("repository:server_catalog:: failed to fetch count of servers %v", err)
	}

	ctr.Page.Total = int(count)

	if ctr.StorageMin != nil || ctr.StorageMax != nil {
		storageQuery := "hdd_size * hdd_count"
		if ctr.StorageMin != nil {
			qry = qry.Where(storageQuery+" >= ?", *ctr.StorageMin)
		}
		if ctr.StorageMax != nil {
			qry = qry.Where(storageQuery+" <= ?", *ctr.StorageMax)
		}
	}

	if len(ctr.RAM) > 0 {
		qry = qry.Where("ram_size IN ?", ctr.RAM)
	}

	if ctr.HDD != nil {
		qry = qry.Where("hdd_type = ?", *ctr.HDD)
	}

	if ctr.Location != nil {
		qry = qry.Where("location = ?", *ctr.Location)
	}

	if err := qry.WithContext(ctx).Limit(ctr.Page.Limit).Offset(ctr.Page.Offset()).Find(&res).Error; err != nil {
		return nil, fmt.Errorf("repository:server_catalog:: failed to fetch  servers %v", err)
	}

	return res, nil
}
