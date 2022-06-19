package envsrc

import (
	"os"

	"github.com/gocombo/config"
	"github.com/gocombo/config/val"
)

type source struct {
	pathToEnvName map[string]string
}

type SourceOpt func(opts *source)

type LoadValOptBuilder struct {
	valuePath string
}

func (b *LoadValOptBuilder) From(envName string) SourceOpt {
	return func(opts *source) {
		opts.pathToEnvName[b.valuePath] = envName
	}
}

func Set(path string) *LoadValOptBuilder {
	return &LoadValOptBuilder{
		valuePath: path,
	}
}

func (s *source) ReadValues(paths []string, optsSetters ...config.ReadValuesOpt) ([]val.Raw, error) {
	opts := &config.ReadValuesOpts{}
	opts.SetOpts(optsSetters...)
	if opts.ChangedSince != nil {
		return nil, nil
	}
	values := make([]val.Raw, 0, len(paths))
	for _, path := range paths {
		envName := s.pathToEnvName[path]
		envVal, ok := os.LookupEnv(envName)
		if !ok {
			continue
		}
		values = append(values, val.Raw{
			Key: path,
			Val: envVal,
		})
	}
	return values, nil
}

func New(opts ...SourceOpt) config.Source {
	s := &source{
		pathToEnvName: make(map[string]string),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
