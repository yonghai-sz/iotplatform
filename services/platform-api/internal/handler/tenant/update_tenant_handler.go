// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package tenant

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"iotplatform/services/platform-api/internal/logic/tenant"
	"iotplatform/services/platform-api/internal/svc"
	"iotplatform/services/platform-api/internal/types"
)

func UpdateTenantHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTenantReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tenant.NewUpdateTenantLogic(r.Context(), svcCtx)
		resp, err := l.UpdateTenant(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
