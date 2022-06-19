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

type source struct {
	fileName          string
	baseDir           string
	ignoreMissingFile bool
	openFile          func(fileName string) (file io.ReadCloser, err error)
}

func newSource() *source {
	return &source{
		baseDir:           "",
		ignoreMissingFile: false,
		openFile: func(fileName string) (file io.ReadCloser, err error) {
			return os.Open(fileName)
		},
	}
}

type SourceOpt func(opts *source)

func WithBaseDir(baseDir string) SourceOpt {
	return func(opts *source) {
		opts.baseDir = baseDir
	}
}

func IgnoreMissingFile() SourceOpt {
	return func(opts *source) {
		opts.ignoreMissingFile = true
	}
}

func getRawValue(key string, source map[string]interface{}) interface{} {
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

func (src *source) ReadValues(keys []string, optsSetters ...config.ReadValuesOpts) ([]val.Raw, error) {
	file, err := src.openFile(path.Join(src.baseDir, src.fileName))
	if err != nil {
		if src.ignoreMissingFile && os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	var rawValuesMap map[string]interface{}
	if err := json.NewDecoder(file).Decode(&rawValuesMap); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}
	result := make([]val.Raw, 0, len(keys))
	for _, key := range keys {
		if v := getRawValue(key, rawValuesMap); v != nil {
			result = append(result, val.Raw{Key: key, Val: v})
		}
	}
	return result, nil
}

func New(fileName string, opts ...SourceOpt) config.Source {
	src := newSource()
	src.fileName = fileName
	for _, opt := range opts {
		opt(src)
	}
	return src
}
