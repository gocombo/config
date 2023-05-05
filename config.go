package config

import (
	"fmt"
	"strings"

	"github.com/gocombo/config/val"
)

type valuesProviderErrors []error

func (v valuesProviderErrors) Error() string {
	result := make([]string, len(v))
	for i, err := range v {
		result[i] = err.Error()
	}
	return "failed building config: " + strings.Join(result, "; ")
}

type valuesProvider struct {
	sources []Source
	errors  valuesProviderErrors
}

// Get returns the value for the given key or false
func (p *valuesProvider) Get(key string) (val.Raw, bool) {
	for i := range p.sources {
		src := p.sources[len(p.sources)-1-i]
		if v, ok := src.GetValue(key); ok {
			return v, true
		}
	}
	return val.Raw{}, false
}

// NotifyError notifies the provider of an error
// that may occur when parsing or is value is missing
func (p *valuesProvider) NotifyError(key string, err error) {
	p.errors = append(p.errors, err)
}

type SourceLoader func() (Source, error)

type Source interface {
	GetValue(key string) (val.Raw, bool)
}

type LoadOpts interface {
	AddSourceLoader(loader SourceLoader)
}

type LoadOpt func(opts LoadOpts)

type configFactory[T any] func(p val.Provider) *T

type loadOpts struct {
	sourceLoaders []SourceLoader
}

func (opts *loadOpts) AddSourceLoader(loader SourceLoader) {
	opts.sourceLoaders = append(opts.sourceLoaders, loader)
}

func Load[T any](factory configFactory[T], optsSetters ...LoadOpt) (*T, error) {
	opts := loadOpts{}
	for _, optSetter := range optsSetters {
		optSetter(&opts)
	}
	if len(opts.sourceLoaders) == 0 {
		return nil, fmt.Errorf("no sources provided")
	}
	sources := make([]Source, len(opts.sourceLoaders))
	for i, loader := range opts.sourceLoaders {
		source, err := loader()
		if err != nil {
			return nil, err
		}
		sources[i] = source
	}

	provider := &valuesProvider{
		sources: sources,
	}

	cfg := factory(provider)
	if provider.errors != nil {
		return nil, provider.errors
	}
	return cfg, nil
}
