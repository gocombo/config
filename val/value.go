package val

type Loader interface {
	Load(path string) (Raw, error)
	NotifyError(path string, err error)
}

func Load[T any](l Loader, key string) T {
	var value T
	raw, err := l.Load(key)
	if err != nil {
		l.NotifyError(key, err)
		return value
	}
	value = raw.Val.(T)
	return value
}
