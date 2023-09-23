package filesrc

import (
	"os"

	"github.com/gocombo/config"
	"github.com/gocombo/config/val"
)

type sourceFileOpt struct {
	filePath      string
	ignoreMissing bool
}

type sourceOpts struct {
	keyToSourceFile map[string]sourceFileOpt
}

type SourceOpt func(opts *sourceOpts)

type LoadValOptBuilder struct {
	valuePath string
}

func (b *LoadValOptBuilder) From(filePath string) SourceOpt {
	return func(opts *sourceOpts) {
		opts.keyToSourceFile[b.valuePath] = sourceFileOpt{
			filePath: filePath,
		}
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
			return load(optSetters...)
		})
	}
}

func load(optSetters ...SourceOpt) (config.Source, error) {
	opts := &sourceOpts{
		keyToSourceFile: map[string]sourceFileOpt{},
	}
	for _, optSetter := range optSetters {
		optSetter(opts)
	}
	src := &source{
		valuesByKey: make(map[string]val.Raw),
	}
	for key, env := range opts.keyToSourceFile {
		data, err := os.ReadFile(env.filePath)
		if err != nil {
			return nil, err
		}
		src.valuesByKey[key] = val.Raw{
			Key: key,
			Val: data,
		}
	}
	return src, nil
}
