package envsrc

import (
	"os"

	"github.com/gocombo/config"
	"github.com/gocombo/config/val"
)

type sourceOpts struct {
	keyToEnvName map[string]string
}

type SourceOpt func(opts *sourceOpts)

type LoadValOptBuilder struct {
	valuePath string
}

func (b *LoadValOptBuilder) From(envName string) SourceOpt {
	return func(opts *sourceOpts) {
		opts.keyToEnvName[b.valuePath] = envName
	}
}

func Set(path string) *LoadValOptBuilder {
	return &LoadValOptBuilder{
		valuePath: path,
	}
}

type source struct {
	valuesByKey map[string]val.Raw
}

func (s *source) GetValue(key string) (val.Raw, bool) {
	rawVal, ok := s.valuesByKey[key]
	if !ok {
		return val.Raw{}, false
	}
	return rawVal, true
}

func Load(optSetters ...SourceOpt) config.LoadOpt {
	return func(opts config.LoadOpts) {
		opts.AddSourceLoader(func() (config.Source, error) {
			return load(optSetters...), nil
		})
	}
}

func load(optSetters ...SourceOpt) config.Source {
	opts := &sourceOpts{
		keyToEnvName: map[string]string{},
	}
	for _, optSetter := range optSetters {
		optSetter(opts)
	}
	src := &source{
		valuesByKey: make(map[string]val.Raw),
	}
	for key, env := range opts.keyToEnvName {
		envVal, ok := os.LookupEnv(env)
		if !ok {
			continue
		}
		src.valuesByKey[key] = val.Raw{
			Key: key,
			Val: envVal,
		}
	}
	return src
}
