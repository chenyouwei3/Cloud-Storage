package logic

import (
	"context"

	"cloud-storage/app/user/api/internal/svc"
	"cloud-storage/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GithubCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// github第三方回调
func NewGithubCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GithubCallbackLogic {
	return &GithubCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GithubCallbackLogic) GithubCallback() (resp *types.TokenResp, err error) {
	// todo: add your logic here and delete this line

	return
}
