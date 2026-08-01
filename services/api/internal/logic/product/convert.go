package product

import (
	"database/sql"
	"time"

	"iotplatform/services/api/internal/types"
	"iotplatform/services/api/model"
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
