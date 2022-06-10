package config

import (
	"fmt"

	"github.com/gocombo/config"
	"github.com/gocombo/config/jsonsrc"
)

type loadOpts struct {
	envName string
}

type LoadOpt func(opts *loadOpts)

func LoadWithEnvName(envName string) LoadOpt {
	return func(opts *loadOpts) {
		opts.envName = envName
	}
}

func defaultLoadOpts() loadOpts {
	return loadOpts{
		envName: "local",
	}
}

func LoadConfig(optSetter ...LoadOpt) *HelloConfig {
	opts := defaultLoadOpts()
	for _, set := range optSetter {
		set(&opts)
	}
	cfg, err := config.Load(
		newConfig,
		config.LoadWithSources(
			jsonsrc.New("default.json"),
			jsonsrc.New(fmt.Sprintf("%s.json", opts.envName)),
			jsonsrc.New(fmt.Sprintf("%s-user.json", opts.envName), jsonsrc.IgnoreMissingFile()),
		),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}
