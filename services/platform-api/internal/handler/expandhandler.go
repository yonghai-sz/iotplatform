// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"iotplatform/services/platform-api/internal/logic"
	"iotplatform/services/platform-api/internal/svc"
	"iotplatform/services/platform-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ExpandHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExpandReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewExpandLogic(r.Context(), svcCtx)
		resp, err := l.Expand(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
