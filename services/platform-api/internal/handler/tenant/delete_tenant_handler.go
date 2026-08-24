// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package tenant

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"iot-zero/services/platform-api/internal/logic/tenant"
	"iot-zero/services/platform-api/internal/svc"
	"iot-zero/services/platform-api/internal/types"
)

func DeleteTenantHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TenantIdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tenant.NewDeleteTenantLogic(r.Context(), svcCtx)
		err := l.DeleteTenant(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
