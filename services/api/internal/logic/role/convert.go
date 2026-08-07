package role

import (
	"database/sql"
	"time"

	"iotplatform/services/api/internal/types"
	"iotplatform/services/api/model"
)

func toRoleInfo(r *model.Role) types.RoleInfo {
	return types.RoleInfo{
		Id:        r.Id,
		RoleKey:   r.RoleKey,
		RoleName:  r.RoleName,
		Enable:    enableToBool(r.Enable),
		TenantId:  r.TenantId,
		CreatedAt: formatNullTime(r.CreatedAt),
		UpdatedAt: formatNullTime(r.UpdatedAt),
	}
}

func formatNullTime(v sql.NullTime) string {
	if !v.Valid {
		return ""
	}
	return v.Time.UTC().Format(time.RFC3339Nano)
}

func boolToEnable(enable bool) string {
	if enable {
		return "Enable"
	}
	return "Disable"
}

func enableToBool(enable string) bool {
	return enable != "Disable"
}
