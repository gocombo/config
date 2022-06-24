package config

import "github.com/gocombo/config/val"

type SourceLoader func() (Source, error)

type Source interface {
	GetValue(keys string) (val.Raw, bool)
}

type LoadOpts interface {
	AddSourceLoader(loader SourceLoader)
}

type LoadOpt func(opts LoadOpts)

type configFactory[T any] func(valLoader val.Provider) *T

func Load[T any](factory configFactory[T], opts ...LoadOpt) (*T, error) {
	return nil, nil
}
