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
	StorageMin *int
	StorageMax *int
	RAM        []int
	HDD        *int
	Location   *string
	Page       *utils.Page `json:"page"`
}

// ListServerResp represents the server information in the response
// @Description Server information in the response
type ListServerResp struct {
	Model    string `json:"model" example:"HP DL120G7Intel G850" description:"Server model name"`
	Ram      string `json:"ram" example:"4GBDDR3" description:"RAM configuration (size and type)"`
	HDD      string `json:"hdd" example:"4x1TBSATA2" description:"Hard disk configuration (count, size and type)"`
	Location string `json:"location" example:"AmsterdamAMS-01" description:"Server location code"`
	Price    string `json:"price" example:"€39.99" description:"Server price with currency symbol"`
}
