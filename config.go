package config

type RawVal struct {
	Val interface{}
}

type RawValSubscription interface {
	subscribe(key string, handler func(val RawVal))
}

type configFactory[T any] func(sub RawValSubscription) *T

type loadOpts struct {
}

type LoadOpt func(opts *loadOpts)

func WithSources(sources ...Source) LoadOpt {
	return nil
}

func Load[T any](newConfig configFactory[T], opts ...LoadOpt) (*T, error) {
	cfg := newConfig(nil)
	return cfg, nil
}
