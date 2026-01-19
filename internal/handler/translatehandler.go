// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/mustang5910/deeplx-translategemma/internal/logic"
	"github.com/mustang5910/deeplx-translategemma/internal/svc"
	"github.com/mustang5910/deeplx-translategemma/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TranslateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewTranslateLogic(r.Context(), svcCtx)
		resp, err := l.Translate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
