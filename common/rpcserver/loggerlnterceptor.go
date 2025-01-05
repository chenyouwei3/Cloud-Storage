package rpcserver

import (
	"cloud-storage/common/errorx"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 参数说明
// ctx 参数说明
// req 请求参数
// info 请求信息的结构体
// handler 处理请求的实际函数
func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Println("err:统一拦截器/", info)
	fmt.Println("req:统一拦截器/", req)

	resp, err = handler(ctx, req)
	if err != nil {
		causeErr := errors.Cause(err)                  //寻找嵌套errors,最根本的错误
		if e, ok := causeErr.(*errorx.CodeError); ok { //如果是自定义错误格式
			logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)
			err = status.Error(codes.Code(e.StatusCode()), e.Error())
		} else {
			logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)
		}
	}
	return resp, err
}
