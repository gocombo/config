package jsonsource

import "github.com/gocombo/config"

type sourceOpts struct {
	ignoreMissingFile bool
}

type SourceOpt func(opts *sourceOpts)

func WithIgnoreMissingFile() SourceOpt {
	return func(opts *sourceOpts) {
		opts.ignoreMissingFile = true
	}
}

func New(filePath string, opts ...SourceOpt) config.Source {
	return nil
}
