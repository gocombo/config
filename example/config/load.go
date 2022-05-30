package config

import (
	"fmt"

	"github.com/gocombo/config"
	"github.com/gocombo/config/jsonsource"
)

type Loader[T any] interface {
	Load() (*T, error)
}

func newConfig(sub config.RawValSubscription) *HelloConfig {
	return &HelloConfig{
		SayHelloTimes: config.MakeVal[int](sub, "sayHelloTimes"),
		Server: &Server{
			Port: config.MakeVal[int](sub, "server/port"),
		},
		Hello: &Hello{
			Message: config.MakeVal[string](sub, "hello/message"),
		},
	}
}

func LoadConfig() *HelloConfig {
	cfg, err := config.Load(newConfig,
		config.WithSources(
			jsonsource.New("config.json"),
		),
	)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	return cfg
}
