package val

type Raw struct {
	Key string
	Val interface{}
}

type Registry interface {
	get(key string, handler func(val Raw) error) bool
}
