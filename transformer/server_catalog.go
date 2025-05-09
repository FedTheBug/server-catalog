package transformer

import (
	"fmt"
	"github.com/server-catalog/internal/dto"
	"github.com/server-catalog/internal/utils"
	"github.com/server-catalog/models"
)

func TransformServerList(servers []models.ServerCatalog) []dto.ListServerResp {
	result := make([]dto.ListServerResp, 0)

	for _, server := range servers {
		var ramType string
		switch server.RamType {
		case utils.RAMTypeDDR3:
			ramType = "DDR3"
		case utils.RAMTypeDDR4:
			ramType = "DDR4"
		default:
			ramType = ""
		}
		ram := fmt.Sprintf("%dGB%s", server.RamSize, ramType)

		hddSize := server.HDDSize
		hddUnit := utils.HDDUnitGB
		if hddSize >= 1024 {
			hddSize = hddSize / 1024
			hddUnit = utils.HDDUnitTB
		}

		var hddType string
		switch server.HDDType {
		case utils.HDDTypeSATA2:
			hddType = utils.HDDSATA2DB
		case utils.HDDTypeSAS:
			hddType = utils.HDDSASDB
		case utils.HDDTypeSSD:
			hddType = utils.HDDSSDDB
		default:
			hddType = ""
		}
		hdd := fmt.Sprintf("%dx%d%s%s", server.HDDCount, hddSize, hddUnit, hddType)

		var price string
		switch server.Currency {
		case 1: // USD
			price = fmt.Sprintf("$%.2f", server.Price)
		case 2: // Euro
			price = fmt.Sprintf("€%.2f", server.Price)
		case 3: // SGD
			price = fmt.Sprintf("S$%.2f", server.Price)
		default:
			price = fmt.Sprintf("%.2f", server.Price)
		}

		result = append(result, dto.ListServerResp{
			Model:    server.Model,
			Ram:      ram,
			HDD:      hdd,
			Location: server.Location,
			Price:    price,
		})

	}
	return result
}
