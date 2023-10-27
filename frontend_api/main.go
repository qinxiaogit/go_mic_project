package main

import (
	"context"
	mgrpc "github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/consul"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	mhttp "github.com/go-micro/plugins/v4/server/http"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qinxiaogit/go_mic_project/frontend_api/config"
	pb "github.com/qinxiaogit/go_mic_project/frontend_api/proto"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"os"
	"time"
)

const (
	name    = "frontend_api"
	version = "1.0.0"

	defaultCurrencyCode = "USD"
	cookieMaxAge        = 60 * 60 * 48
	cookiePrefix        = "hipstershop_"
	cookiesSessionID    = cookiePrefix + "session-id"
	cookieCurrency      = cookiePrefix + "currency"
)

var (
	whitelistedCurrencies = map[string]bool{
		"USD": true,
		"EUR": true,
		"CAD": true,
		"JPY": true,
		"GBP": true,
		"TRY": true,
	}
)

type ctxKeySessionID struct {
}
type frontendApiServer struct {
	shippingService pb.ShippingService
	cacheService    pb.CacheService
}

func main() {
	if err := config.Load(); err != nil {
		logger.Fatal(err)
	}
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	// Create service
	srv := micro.NewService(
		micro.Server(mhttp.NewServer()),
		micro.Client(mgrpc.NewClient()),
		micro.Registry(consulReg),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
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
	}
	srv.Init(opts...)

	// Register handler
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout

	cfg, client := config.Get(), srv.Client()
	svc := &frontendApiServer{
		shippingService: pb.NewShippingService(cfg.ShippingService, client),
		cacheService:    pb.NewCacheService(cfg.CacheService, client),
	}

	r := mux.NewRouter()

	r.HandleFunc("/quote", svc.GetQuote).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/order", svc.ShipOrder).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/set", svc.setCache).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/get", svc.getCache).Methods(http.MethodGet, http.MethodHead)
	r.Handle("/metrics", promhttp.Handler())

	var handler http.Handler = r
	handler = &logHandler{log: log, next: handler} // add logging
	handler = ensureSessionID(handler)             // add session ID
	// handler = tracing(handler)                     // add opentelemetry instrumentation
	r.Use(otelmux.Middleware(name))
	r.Use(tracingContextWrapper)
	if err := micro.RegisterHandler(srv.Server(), handler); err != nil {
		logger.Fatal(err)
	}

	logger.Infof("starting server on %s", config.Address())
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
