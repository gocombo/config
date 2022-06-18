package jsonsrc

import "github.com/gocombo/config"

type sourceOpts struct {
	ignoreMissingFile bool
}

type SourceOpt func(opts *sourceOpts)

func IgnoreMissingFile() SourceOpt {
	return func(opts *sourceOpts) {
		opts.ignoreMissingFile = true
	}
}

func New(fileName string, opts ...SourceOpt) config.Source {
	return nil
}
