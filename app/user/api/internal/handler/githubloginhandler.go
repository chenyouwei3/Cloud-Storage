package handler

import (
	"net/http"

	"cloud-storage/app/user/api/internal/logic"
	"cloud-storage/app/user/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// github第三方登陆
func GithubLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGithubLoginLogic(r.Context(), svcCtx)
		err := l.GithubLogin()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
