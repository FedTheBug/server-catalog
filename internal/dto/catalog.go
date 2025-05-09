package dto

import (
	"github.com/server-catalog/internal/utils"
	"mime/multipart"
)

// UploadCatalogCtr ...
type UploadCatalogCtr struct {
	File multipart.File
}

// ListServersCtr ...
type ListServersCtr struct {
	Page *utils.Page `json:"page"`
}

type ListServerResp struct {
	Model    string `json:"model"`
	Ram      string `json:"ram"`
	HDD      string `json:"hdd"`
	Location string `json:"location"`
	Price    string `json:"price"`
}
