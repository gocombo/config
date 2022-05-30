package config

type Value[T any] func() T

func MakeVal[T any](sub RawValSubscription, key string) Value[T] {
	var result T
	return func() T {
		return result
	}
}
