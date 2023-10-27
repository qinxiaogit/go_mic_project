package main

import (
	"context"
	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"strings"
	"time"

	"github.com/qinxiaogit/go_mic_project/payment_service/config"
	"github.com/qinxiaogit/go_mic_project/payment_service/handler"
	pb "github.com/qinxiaogit/go_mic_project/payment_service/proto"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

var (
	name    = "payment_service"
	version = "1.0.0"
)

func main() {
	if err := config.Load(); err != nil {
		logger.Fatal(err)
	}
	//create service
	srv := micro.NewService(
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
	)

	opts := []micro.Option{
		micro.Name(name),
		micro.Version(version),
		micro.Address(config.Address()),
	}

	// Register handler
	if cfg := config.Tracing(); cfg.Enable {
		tp, err := newTracerProvider(name, srv.Server().Options().Id, cfg.Jaeger.Url)
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
	// Register handler
	if err := pb.RegisterPaymentServiceHandler(srv.Server(), new(handler.PaymentService)); err != nil {
		logger.Fatal(err)
	}
	if err := pb.RegisterHealthHandler(srv.Server(), new(handler.Health)); err != nil {

	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

func newTracerProvider(servername, serverId, url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(servername),
			semconv.ServiceVersionKey.String(version),
			semconv.ServiceInstanceIDKey.String(serverId),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
