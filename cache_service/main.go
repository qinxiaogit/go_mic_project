package main

import (
	"context"
	"github.com/go-micro/plugins/v4/registry/consul"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	"github.com/qinxiaogit/go_mic_project/cache_service/config"
	"github.com/qinxiaogit/go_mic_project/cache_service/handler"
	pb "github.com/qinxiaogit/go_mic_project/cache_service/proto"
	"go-micro.dev/v4"
	"go-micro.dev/v4/cache"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"strings"
	"time"

	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
)

var (
	name    = "cache_service"
	version = "1.0.0"
)

func main() {
	if err := config.Load(); err != nil {
		logger.Fatal(err)
	}
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	// Create service
	srv := micro.NewService(
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
		micro.Registry(consulReg),
	)

	opts := []micro.Option{
		micro.Name(name),
		micro.Version(version),
		micro.Address(config.Address()),
	}

	if cfg := config.Tracing(); cfg.Enable {
		tp, err := newTracerProvider(name, srv.Server().Options().Id, cfg.Jaeger.URL)
		if err != nil {
			logger.Fatal(err)
		}
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				logger.Fatal(err)
			}
		}()

		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		traceOpts := []opentelemetry.Option{
			opentelemetry.WithHandleFilter(func(ctx context.Context, request server.Request) bool {
				if e := request.Endpoint(); strings.HasPrefix(e, "Health.") {
					return true
				}
				return false
			}),
		}
		opts = append(opts, micro.WrapHandler(opentelemetry.NewHandlerWrapper(traceOpts...)))
	}
	srv.Init(opts...)
	//Register handler
	cacheOption := []cache.Option{
		func(o *cache.Options) {
			o.Address = "127.0.0.1:6379"
		},
	}
	// Register handler
	if err := pb.RegisterCacheHandler(srv.Server(), handler.NewCache(cacheOption...)); err != nil {
		logger.Fatal(err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
