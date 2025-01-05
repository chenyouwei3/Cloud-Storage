package handler

import (
	"net/http"

	"cloud-storage/app/user/api/internal/logic"
	"cloud-storage/app/user/api/internal/svc"
	"cloud-storage/app/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 发验证码
func SendcodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterByPhoneRep
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSendcodeLogic(r.Context(), svcCtx)
		resp, err := l.Sendcode(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
