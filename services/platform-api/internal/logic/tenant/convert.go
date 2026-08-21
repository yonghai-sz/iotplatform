package tenant

import (
	"database/sql"
	"time"

	"iotplatform/services/platform-api/internal/types"
	"iotplatform/services/platform-api/model"
)

func toTenantInfo(t *model.Tenants) types.TenantInfo {
	return types.TenantInfo{
		Id:         t.Id,
		TenantName: t.TenantName,
		Email:      nullString(t.Email),
		CreatedAt:  formatNullTime(t.CreatedAt),
		UpdatedAt:  formatNullTime(t.UpdatedAt),
	}
}

func nullString(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

func formatNullTime(v sql.NullTime) string {
	if !v.Valid {
		return ""
	}
	return v.Time.UTC().Format(time.RFC3339Nano)
}

func emailToNullString(email string) sql.NullString {
	if email == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: email, Valid: true}
}
