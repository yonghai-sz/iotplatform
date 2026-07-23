// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package tenant

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"iotplatform/services/api/internal/logic/tenant"
	"iotplatform/services/api/internal/svc"
	"iotplatform/services/api/internal/types"
)

func ListTenantHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListTenantReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tenant.NewListTenantLogic(r.Context(), svcCtx)
		resp, err := l.ListTenant(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
