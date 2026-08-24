package product

import (
	"database/sql"
	"time"

	"iot-zero/services/platform-api/internal/types"
	"iot-zero/services/platform-api/model"
)

func toProductInfo(p *model.Product) types.ProductInfo {
	return types.ProductInfo{
		Id:          p.Id,
		ProductCode: p.ProductCode,
		ProductName: p.ProductName,
		CreatedAt:   formatNullTime(p.CreatedAt),
		UpdatedAt:   formatNullTime(p.UpdatedAt),
	}
}

func formatNullTime(v sql.NullTime) string {
	if !v.Valid {
		return ""
	}
	return v.Time.UTC().Format(time.RFC3339Nano)
}
