package usecase

import (
	"bytes"
	"context"
	"fmt"
	"github.com/server-catalog/internal/dto"
	"github.com/server-catalog/internal/utils"
	"github.com/server-catalog/models"
	"github.com/server-catalog/repository"
	"github.com/server-catalog/transformer"
	"github.com/xuri/excelize/v2"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type ServerCatalog struct {
	SCRepo repository.CatalogRepository
}

func New(scr repository.CatalogRepository) CatalogUseCase {
	return &ServerCatalog{SCRepo: scr}
}

func (sc *ServerCatalog) UploadCatalog(ctx context.Context, ctr *dto.UploadCatalogCtr) error {
	data, err := io.ReadAll(ctr.File)
	if err != nil {
		return fmt.Errorf("usecase:server_catalog:failed to read file")
	}
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("usecase:server_catalog:invalid XLSX file")
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return fmt.Errorf("usecase:server_catalog:no sheet found")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) < 2 {
		return fmt.Errorf("usecase:server_catalog:no data in the sheet")
	}

	if len(rows)-1 > 1000 {
		return fmt.Errorf("usecase:server_catalog:maximum number of rows exceeded (1000)")
	}

	// Validate header
	expected := []string{"Model", "RAM", "HDD", "Location", "Price"}
	header := rows[0]
	for i, col := range expected {
		if i >= len(header) || strings.TrimSpace(header[i]) != col {
			return fmt.Errorf("usecase:server_catalog:invalid XLSX columns")
		}
	}

	var inserted int
	catalogs := make([]models.ServerCatalog, 0)
	for idx, row := range rows[1:] {
		if len(row) < 5 {
			return fmt.Errorf("usecase:server_catalog:rows too short")
		}

		model := strings.TrimSpace(row[0])
		if model == "" {
			return fmt.Errorf("usecase:server_catalog:model is required at %d", idx+2)
		}

		ramSize, ramType, err := parseRAM(row[1])
		if err != nil {
			return fmt.Errorf("usecase:server_catalog:Invalid RAM at row %d: %s\n", idx+2, row[1])
		}

		hddCount, hddSize, hddType, err := parseHDD(row[2])
		if err != nil {
			return fmt.Errorf("usecase:server_catalog:Invalid HDD at row %d: %s", idx+2, row[2])
		}

		location := strings.TrimSpace(row[3])
		if location == "" {
			return fmt.Errorf("usecase:server_catalog:Location is required at row %d", idx+2)
		}

		price, currencySymbol, err := parsePrice(row[4])
		if err != nil {
			return fmt.Errorf("usecase:server_catalog:Invalid Price at row %d: %s", idx+2, row[4])
		}

		ramTypeID, err := utils.GetRAMTypeID(ramType)
		if err != nil {
			return err
		}

		hddTypeID, err := utils.GetHDDTypeID(hddType)
		if err != nil {
			return err
		}

		currencyID, err := utils.GetCurrencyID(currencySymbol)
		if err != nil {
			return err
		}

		catalog := models.ServerCatalog{
			Model:    model,
			RamSize:  ramSize,
			RamType:  ramTypeID,
			HDDSize:  hddSize,
			HDDCount: hddCount,
			HDDType:  hddTypeID,
			Location: location,
			Price:    price,
			Currency: currencyID,
		}

		catalogs = append(catalogs, catalog)
		inserted++
	}

	if err := sc.SCRepo.Upload(ctx, catalogs); err != nil {
		return fmt.Errorf("usecase:server_catalog:: failed to upload %v", utils.ErrUploadFailed)
	}
	return nil
}

func parseRAM(ram string) (int, string, error) {
	ram = strings.ToUpper(strings.TrimSpace(ram))
	re := regexp.MustCompile(`^(\d+)\s*GB\s*(\w+)$`)
	matches := re.FindStringSubmatch(ram)
	if len(matches) != 3 {
		return 0, "", fmt.Errorf("invalid RAM format")
	}
	size, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", err
	}
	typ := matches[2]
	return size, typ, err
}

func parseHDD(hdd string) (int, int, string, error) {
	hdd = strings.TrimSpace(hdd)
	re := regexp.MustCompile(`(?i)^(\d+)x(\d+)(TB|GB)([A-Z0-9]+)$`)
	matches := re.FindStringSubmatch(hdd)
	if len(matches) != 5 {
		return 0, 0, "", fmt.Errorf("invalid HDD format: %q", hdd)
	}

	count, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, "", err
	}

	size, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, "", err
	}

	unit := strings.ToUpper(matches[3])
	typ := strings.ToUpper(matches[4])

	// Convert to GB
	sizeGB := utils.ConvertToGB(size, unit)

	return count, sizeGB, typ, nil
}

func parsePrice(price string) (float64, string, error) {
	price = strings.TrimSpace(price)

	re := regexp.MustCompile(`^([^\d]+)\s*([\d\.]+)$`)
	matches := re.FindStringSubmatch(price)
	if len(matches) != 3 {
		return 0, "", fmt.Errorf("invalid price format: %s", price)
	}

	currency := strings.TrimSpace(matches[1])

	if currency == "S$" {
		currency = utils.CurrencySymbolSGD
	}

	amount, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid amount format: %s", matches[2])
	}

	return amount, currency, nil
}

func (sc *ServerCatalog) GetLocations(ctx context.Context) ([]string, error) {
	locs, err := sc.SCRepo.GetLocations(ctx)
	if err != nil {
		return nil, err
	}

	return locs, nil
}

func (sc *ServerCatalog) GetHDDTypes(ctx context.Context) ([]string, error) {
	hddTypes, err := sc.SCRepo.GetHDDTypes(ctx)
	if err != nil {
		return nil, err
	}

	return hddTypes, nil
}

func (sc *ServerCatalog) GetListOfServers(ctx context.Context, ctr *dto.ListServersCtr) ([]dto.ListServerResp, error) {
	result, err := sc.SCRepo.GetServers(ctx, ctr)
	if err != nil {
		return nil, fmt.Errorf("usecase:server_catalog:: failed to get list of servers %v", err)
	}

	if len(result) < 1 {
		return nil, utils.ErrServerNotFound
	}

	transformedList := transformer.TransformServerList(result)

	return transformedList, nil
}
