package dto

import "mime/multipart"

// UploadCatalogCtr ...
type UploadCatalogCtr struct {
	File multipart.File
}
