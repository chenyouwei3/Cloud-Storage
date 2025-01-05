package main

import (
	"cloud-storage/common/errorx"
	"cloud-storage/common/logs/zapx"
	"cloud-storage/common/rpcserver"
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"cloud-storage/app/user/rpc/internal/config"
	"cloud-storage/app/user/rpc/internal/server"
	"cloud-storage/app/user/rpc/internal/svc"
	"cloud-storage/app/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	//解析程序启动时传入的命令行参数
	flag.Parse()
	// 关闭任何之前的日志记录器资源
	logx.Close()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	//添加单一拦截器
	s.AddUnaryInterceptors(rpcserver.LoggerInterceptor)
	//自定义错误
	//httpx.SetErrorHandlerCtx 这个函数确实是规定了一个统一的错误处理机制，确保了所有错误都必须通过这个自定义的方式返回给客户端
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx.CodeError:
			return http.StatusOK, e.Data()
		default:
			return http.StatusInternalServerError, nil
		}
	})
	writer, err := zapx.NewZapWriter()
	logx.Must(err)
	logx.SetWriter(writer)
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
