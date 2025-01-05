package handler

import (
	"net/http"

	"cloud-storage/app/user/api/internal/logic"
	"cloud-storage/app/user/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// github第三方回调
func GithubCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGithubCallbackLogic(r.Context(), svcCtx)
		resp, err := l.GithubCallback()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
