package envsrc

import "github.com/gocombo/config"

type sourceOpts struct {
	pathToEnvName map[string]string
}

type SourceOpt func(opts *sourceOpts)

type LoadValOptBuilder struct {
	valuePath string
}

func (b *LoadValOptBuilder) FromEnv(envName string) SourceOpt {
	return func(opts *sourceOpts) {
		opts.pathToEnvName[b.valuePath] = envName
	}
}

func SetValue(path string) *LoadValOptBuilder {
	return &LoadValOptBuilder{
		valuePath: path,
	}
}

func New(opts ...SourceOpt) config.Source {
	return nil
}
