package val

import "fmt"

type Provider interface {
	Get(key string) (Raw, error)
	NotifyError(key string, err error)
}

func Get[T any](l Provider, key string) T {
	var value T
	raw, err := l.Get(key)
	if err != nil {
		l.NotifyError(key, err)
		return value
	}
	value, ok := raw.Val.(T)
	if !ok {
		l.NotifyError(key, fmt.Errorf("value not a %T: %s", value, key))
	}
	return value
}
