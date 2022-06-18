package val

type Raw struct {
	Val interface{}
}

type Registry interface {
	get(key string, handler func(val Raw) error) bool
}
