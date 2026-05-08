// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"gozero-example-iot/services/mqtt-auth-api/internal/logic"
	"gozero-example-iot/services/mqtt-auth-api/internal/svc"
	"gozero-example-iot/services/mqtt-auth-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAuthLogic(r.Context(), svcCtx)
		resp, err := l.Auth(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
