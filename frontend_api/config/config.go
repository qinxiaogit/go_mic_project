package config

import (
	"github.com/go-micro/plugins/v4/config/source/consul"
	"github.com/pkg/errors"
	"go-micro.dev/v4/config"
)

type Config struct {
	Address               string
	Tracing               TracingConfig
	AdService             string
	CartService           string
	CheckoutService       string
	CurrencyService       string
	ProductCatalogService string
	RecommendationService string
	ShippingService       string
	CacheService          string
}

type TracingConfig struct {
	Enable bool
	Jaeger JaegerConfig
}

type JaegerConfig struct {
	URL string
}

var cfg *Config = &Config{
	Address:               ":8090",
	AdService:             "adservice",
	CartService:           "cartservice",
	CheckoutService:       "checkoutservice",
	CurrencyService:       "currencyservice",
	ProductCatalogService: "productcatalogservice",
	RecommendationService: "recommendationservice",
	ShippingService:       "shiping_service",
	CacheService:          "cache_service",
}

func Get() Config {
	return *cfg
}

func Address() string {
	return cfg.Address
}

func Tracing() TracingConfig {
	return cfg.Tracing
}

func Load() error {
	consulSource := consul.NewSource(consul.WithAddress("127.0.0.1:8500"), consul.WithPrefix(""))
	configor, err := config.NewConfig(config.WithSource(consulSource))
	if err != nil {
		return errors.Wrap(err, "configor.New")
	}
	if err := configor.Load(); err != nil {
		return errors.Wrap(err, "configor.Load")
	}
	if err := configor.Scan(cfg); err != nil {
		return errors.Wrap(err, "configor.Scan")
	}
	//go func() {
	//	fmt.Println("加载配置文件")
	//	watch, err := configor.Watch("")
	//	fmt.Println(1, err)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	v, err := watch.Next()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(2)
	//
	//	if err := v.Scan(cfg); err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println("配置文件加载完毕")
	//}()
	return nil
}
