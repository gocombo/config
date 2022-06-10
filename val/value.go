package val

type Loader interface {
	Load(path string) (Raw, error)
	NotifyError(path string, err error)
}

func Load[T any](l Loader, key string) T {
	var value T
	return value
}
