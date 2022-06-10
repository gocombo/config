package config

import (
	"time"

	"github.com/gocombo/config/val"
)

type ReadValuesOpts struct {
	ChangedSince *time.Time
}

func (opts *ReadValuesOpts) SetOpts(optsSetters ...ReadValuesOpt) {
	for _, optsSetter := range optsSetters {
		optsSetter(opts)
	}
}

type ReadValuesOpt func(opts *ReadValuesOpts)

type Source interface {
	ReadValues(keys []string, optsSetters ...ReadValuesOpts) ([]val.Raw, error)
}
