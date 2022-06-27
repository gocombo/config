package val

import "fmt"

type Raw struct {
	Key string
	Val interface{}
}

type Provider interface {
	// Get returns the value for the given key or false
	Get(key string) (Raw, bool)

	// NotifyError notifies the provider of an error
	// that may occur when parsing or is value is missing
	NotifyError(key string, err error)
}

func Define[T any](l Provider, key string) T {
	var value T
	raw, ok := l.Get(key)
	if !ok {
		l.NotifyError(key, fmt.Errorf("value %s not found", key))
		return value
	}
	value, ok = raw.Val.(T)
	if !ok {
		l.NotifyError(key, fmt.Errorf("value not a %T: %s", value, key))
	}
	return value
}
