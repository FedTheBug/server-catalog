package utils

import (
	"encoding/json"
	"fmt"
	_ "github.com/davecgh/go-spew/spew"
	"io"
	"strings"
)

// HDD Types
const (
	HDDTypeSATA2 = 1
	HDDTypeSAS   = 2
	HDDTypeSSD   = 3
)

// RAM Types
const (
	RAMTypeDDR3 = 1
	RAMTypeDDR4 = 2
)

// Currency Types
const (
	CurrencyUSD  = 1
	CurrencyEuro = 2
	CurrencySGD  = 3
)

// Currency Symbols
const (
	CurrencySymbolUSD  = "$"
	CurrencySymbolEuro = "€"
	CurrencySymbolSGD  = "S$"
)

const (
	StorageUnitGB = 1
	StorageUnitTB = 1024 // 1 TB = 1024 GB

	HDDUnitGB = "GB"
	HDDUnitTB = "TB"
)

// GetHDDTypeID returns the HDD type ID based on the parsed type string.
func GetHDDTypeID(hddType string) (int, error) {
	hddType = strings.ToUpper(hddType)
	switch hddType {
	case "SATA2":
		return HDDTypeSATA2, nil
	case "SAS":
		return HDDTypeSAS, nil
	case "SSD":
		return HDDTypeSSD, nil
	default:
		return 0, fmt.Errorf("unknown HDD type: %s", hddType)
	}
}

// GetRAMTypeID returns the RAM type ID based on the parsed type string.
func GetRAMTypeID(ramType string) (int, error) {
	ramType = strings.ToUpper(ramType)
	switch ramType {
	case "DDR3":
		return RAMTypeDDR3, nil
	case "DDR4":
		return RAMTypeDDR4, nil
	default:
		return 0, fmt.Errorf("unknown RAM type: %s", ramType)
	}
}

// ConvertToGB return the size in Gigabytes
func ConvertToGB(size int, unit string) int {
	switch strings.ToUpper(unit) {
	case "TB":
		return size * StorageUnitTB
	case "GB":
		return size * StorageUnitGB
	default:
		return size
	}
}

// GetCurrencyID returns the currency ID based on the parsed currency symbol.
func GetCurrencyID(currencySymbol string) (int, error) {
	currencySymbol = strings.TrimSpace(currencySymbol)

	switch currencySymbol {
	case CurrencySymbolUSD:
		return CurrencyUSD, nil
	case CurrencySymbolEuro:
		return CurrencyEuro, nil
	case CurrencySymbolSGD:
		return CurrencySGD, nil
	default:
		return 0, fmt.Errorf("unknown currency symbol: %s", currencySymbol)
	}
}

// ParseJSON parses the JSON response body into the provided interface
func ParseJSON(body io.Reader, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}
