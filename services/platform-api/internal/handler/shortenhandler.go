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

func ShortenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShortenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShortenLogic(r.Context(), svcCtx)
		resp, err := l.Shorten(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
