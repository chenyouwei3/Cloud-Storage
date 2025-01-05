package handler

import (
	"net/http"

	"cloud-storage/app/user/api/internal/logic"
	"cloud-storage/app/user/api/internal/svc"
	"cloud-storage/app/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户信息
func userInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.UserInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
