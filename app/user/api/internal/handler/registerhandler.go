package handler

import (
	"net/http"

	"cloud-storage/app/user/api/internal/logic"
	"cloud-storage/app/user/api/internal/svc"
	"cloud-storage/app/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户验证码形式注册/登陆
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}