package config

import "github.com/gocombo/config/val"

type loadOpts struct {
}

type LoadOpt func(opts *loadOpts)

type configFactory[T any] func(valLoader val.Loader) *T

func LoadWithSources(sources ...Source) LoadOpt {
	return func(opts *loadOpts) {
	}
}

func Load[T any](factory configFactory[T], opts ...LoadOpt) (*T, error) {
	return nil, nil
}
