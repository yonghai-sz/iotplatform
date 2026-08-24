package user

import (
	"database/sql"
	"time"

	"iot-zero/services/platform-api/internal/types"
	"iot-zero/services/platform-api/model"
)

func toUserInfo(u *model.User) types.UserInfo {
	return types.UserInfo{
		Id:        u.Id,
		Username:  u.Username,
		Enable:    enableToBool(u.Enable),
		RoleId:    u.RoleId,
		TenantId:  u.TenantId,
		CreatedAt: formatNullTime(u.CreatedAt),
		UpdatedAt: formatNullTime(u.UpdatedAt),
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
