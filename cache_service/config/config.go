package config

import (
	"fmt"
	"github.com/go-micro/plugins/v4/config/source/consul"
	"github.com/pkg/errors"
	"go-micro.dev/v4/config"
)

type Config struct {
	Port    int
	Tracing TracingConfig
}

type TracingConfig struct {
	Enable bool
	Jaeger JaegerConfig
}

type JaegerConfig struct {
	URL string
}

var cfg *Config = &Config{
	Port: 50057,
}

func Address() string {
	return fmt.Sprintf(":%d", cfg.Port)
}

func Tracing() TracingConfig {
	return cfg.Tracing
}

func Load() error {
	consulSource := consul.NewSource(consul.WithAddress("127.0.0.1:8500"), consul.WithPrefix("tracing"))

	configor, err := config.NewConfig(config.WithSource(consulSource))
	if err != nil {
		return errors.Wrap(err, "configor.New")
	}
	if err := configor.Load(); err != nil {
		return errors.Wrap(err, "configor.Load")
	}
	fmt.Println(string(configor.Bytes()))
	if err := configor.Scan(cfg); err != nil {
		return errors.Wrap(err, "configor.Scan")
	}
	fmt.Println(cfg)
	return nil
}
