package utils

const (
	// Frontend display values
	HDDSATAFE = "SATA"
	HDDSASFE  = "SAS"
	HDDSSDFE  = "SSD"

	// Database values
	HDDSATA2DB = "SATA2"
	HDDSASDB   = "SAS"
	HDDSSDDB   = "SSD"
)

// Mapping from frontend to database values
var HDDTypeMapping = map[string]string{
	HDDSATAFE: HDDSATA2DB,
	HDDSASFE:  HDDSASDB,
	HDDSSDFE:  HDDSSDDB,
}

// Reverse mapping from database to frontend values
var HDDTypeReverseMapping = map[string]string{
	HDDSATA2DB: HDDSATAFE,
	HDDSASDB:   HDDSASFE,
	HDDSSDDB:   HDDSSDFE,
}
