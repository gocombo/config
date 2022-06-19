package config

import (
	"fmt"

	"github.com/gocombo/config"
	"github.com/gocombo/config/envsrc"
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
			// Order of data sources makes a difference
			// Last one has a priority and will "override" values
			// of a previous one, if such values are available in a source.

			// default.json defines base config with initial (default) values
			jsonsrc.New("default.json", jsonsrc.WithBaseDir("config")),

			// environment specific file allows overriding defaults
			jsonsrc.New(fmt.Sprintf("%s.json", opts.envName), jsonsrc.WithBaseDir("config")),

			// <user> configs can be used to let devs override values locally without committing
			jsonsrc.New(fmt.Sprintf("%s-user.json", opts.envName),
				jsonsrc.WithBaseDir("config"),
				jsonsrc.IgnoreMissingFile(),
			),

			// Allow overriding some values via environment variables
			envsrc.New(
				envsrc.Set("server/port").From("PORT"),
			),
		),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}
