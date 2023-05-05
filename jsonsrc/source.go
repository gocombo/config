package jsonsrc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/gocombo/config"
	"github.com/gocombo/config/val"
)

type loadOpts struct {
	baseDir           string
	ignoreMissingFile bool
	openFile          func(fileName string) (file io.ReadCloser, err error)
}

func defaultLoadOpts() loadOpts {
	return loadOpts{
		baseDir:           "",
		ignoreMissingFile: false,
		openFile: func(fileName string) (file io.ReadCloser, err error) {
			return os.Open(fileName)
		},
	}
}

func (o *loadOpts) set(optSetter []LoadOpt) {
	for _, opt := range optSetter {
		opt(o)
	}
}

type LoadOpt func(opts *loadOpts)

func WithBaseDir(baseDir string) LoadOpt {
	return func(opts *loadOpts) {
		opts.baseDir = baseDir
	}
}

func IgnoreMissingFile() LoadOpt {
	return func(opts *loadOpts) {
		opts.ignoreMissingFile = true
	}
}

func getRawValue(key string, source map[string]interface{}) interface{} {
	if source == nil {
		return nil
	}
	firstSeparatorIndex := strings.Index(key, "/")
	if firstSeparatorIndex >= 0 {
		parentKey := key[:firstSeparatorIndex]
		nestedKey := key[firstSeparatorIndex+1:]
		if nestedSource, ok := source[parentKey].(map[string]interface{}); ok {
			return getRawValue(
				nestedKey,
				nestedSource,
			)
		}
	}
	if v, ok := source[key]; ok {
		return v
	}
	return nil
}

type source struct {
	rawValues map[string]interface{}
}

// TODO: Null value support
func (src *source) GetValue(key string) (val.Raw, bool) {
	if v := getRawValue(key, src.rawValues); v != nil {
		return val.Raw{Key: key, Val: v}, true
	}
	return val.Raw{}, false
}

func Load(fileName string, optSetter ...LoadOpt) config.LoadOpt {
	return func(opts config.LoadOpts) {
		opts.AddSourceLoader(func() (config.Source, error) {
			return load(fileName, optSetter...)
		})
	}
}

func load(fileName string, optSetter ...LoadOpt) (config.Source, error) {
	opts := defaultLoadOpts()
	opts.set(optSetter)

	file, err := opts.openFile(path.Join(opts.baseDir, fileName))
	if err != nil {
		if opts.ignoreMissingFile && os.IsNotExist(err) {
			return &source{}, nil
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	src := source{
		rawValues: map[string]interface{}{},
	}
	if err := json.NewDecoder(file).Decode(&src.rawValues); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}
	return &src, nil
}
