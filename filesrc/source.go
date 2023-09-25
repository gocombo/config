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

type FromOpt func(s *sourceFileOpt)

func IgnoreMissing() FromOpt {
	return func(s *sourceFileOpt) {
		s.ignoreMissing = true
	}
}

// From defines filePath to set value from
func (b *LoadValOptBuilder) From(filePath string, pathOpts ...FromOpt) SourceOpt {
	return func(opts *sourceOpts) {
		s := sourceFileOpt{
			filePath: filePath,
		}
		for _, pathOpt := range pathOpts {
			pathOpt(&s)
		}
		opts.keyToSourceFile[b.valuePath] = s
	}
}

// Set config value using From
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
		isMissing := os.IsNotExist(err)
		if err != nil && !(env.ignoreMissing && isMissing) {
			return nil, err
		}
		if isMissing {
			continue
		}
		src.valuesByKey[key] = val.Raw{
			Key: key,
			Val: data,
		}
	}
	return src, nil
}
