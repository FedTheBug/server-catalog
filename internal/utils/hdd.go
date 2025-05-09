package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// Database values
	HDDSATA2DB = "SATA2"
	HDDSASDB   = "SAS"
	HDDSSDDB   = "SSD"
)

func ParseStorageToGB(storage string) (int, error) {
	storage = strings.TrimSpace(storage)
	if storage == "0" {
		return 0, nil
	}

	// GB
	if strings.HasSuffix(storage, "GB") {
		numStr := strings.TrimSuffix(storage, "GB")
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, fmt.Errorf("invalid GB format: %s", storage)
		}
		return num, nil
	}

	// TB
	if strings.HasSuffix(storage, "TB") {
		numStr := strings.TrimSuffix(storage, "TB")
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, fmt.Errorf("invalid TB format: %s", storage)
		}
		return num * 1024, nil // Convert TB to GB (1 TB = 1024 GB)
	}

	return 0, fmt.Errorf("invalid storage format: %s", storage)
}
