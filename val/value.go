package val

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ErrConvertFailed struct {
	message        string
	source         interface{}
	targetTypeName string
}

func (e ErrConvertFailed) Error() string {
	return fmt.Sprintf("failed to convert %[1]v{%[1]T} to %v: %s", e.source, e.targetTypeName, e.message)
}

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

type typeConverter map[string]func(source interface{}, target reflect.Value) error

var supportedConverters = typeConverter{
	"string": func(val interface{}, target reflect.Value) error {
		if _, ok := val.(string); !ok {
			return fmt.Errorf("not a string")
		}
		targetVal := reflect.ValueOf(val).Convert(target.Type())
		target.Set(targetVal)
		return nil
	},
	"[]string": func(val interface{}, target reflect.Value) error {
		var strSlice []string
		var err error
		switch actualSliceVal := val.(type) {
		case string:
			strSlice = strings.FieldsFunc(strings.Replace(actualSliceVal, " ", "", -1), func(c rune) bool {
				return c == ','
			})
		case []string:
			strSlice = actualSliceVal
		case []interface{}: // By default json.Unmarshal will use []interface{} for string arrays
			strSlice = make([]string, len(actualSliceVal))
		parseActualVal:
			for i, val := range actualSliceVal {
				strVal, ok := val.(string)
				if !ok {
					err = fmt.Errorf("expected []string but got %v(%[1]T)", actualSliceVal)
					break parseActualVal
				}
				strSlice[i] = strVal
			}
		default:
			err = fmt.Errorf("expected []string but got %v(%[1]T)", val)
		}
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(strSlice))
		return nil
	},
	"int": func(val interface{}, target reflect.Value) error {
		var intVal int
		var err error
		switch actualVal := val.(type) {
		case int:
			intVal = actualVal
		// case float32:
		// 	intVal = int(actualVal)
		// case float64:
		// 	intVal = int(actualVal)
		case string:
			intVal, err = strconv.Atoi(actualVal)
		default:
			err = errors.New("unexpected int type")
		}
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(intVal))
		return nil
	},
	"int64": func(val interface{}, target reflect.Value) error {
		var intVal int64
		var err error
		switch actualVal := val.(type) {
		case int64:
			intVal = actualVal
		case int32:
			intVal = int64(actualVal)
		case int:
			intVal = int64(actualVal)
		case string:
			intVal, err = strconv.ParseInt(actualVal, 10, 64)
		default:
			err = errors.New("unexpected int64 type")
		}
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(intVal))
		return nil
	},
	"float64": func(val interface{}, target reflect.Value) error {
		var floatVal float64
		var err error
		switch actualVal := val.(type) {
		case int:
			floatVal = float64(actualVal)
		case float32:
			floatVal = float64(actualVal)
		case float64:
			floatVal = actualVal
		case string:
			floatVal, err = strconv.ParseFloat(actualVal, 64)
		default:
			err = errors.New("unexpected float64 type")
		}
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(floatVal))
		return nil
	},
	"bool": func(val interface{}, target reflect.Value) error {
		var boolVal bool
		var err error
		switch newVal := val.(type) {
		case bool:
			boolVal = newVal
		case string:
			boolVal, err = strconv.ParseBool(newVal)
		}
		if err != nil {
			return fmt.Errorf("Expected bool value but got: %v(%[1]T)", val)
		}
		target.Set(reflect.ValueOf(boolVal))
		return nil
	},
}

func (c typeConverter) convert(source interface{}, target reflect.Value) error {
	kind := target.Kind()
	typeName := kind.String()

	if kind == reflect.Slice {
		typeName = "[]" + target.Type().Elem().Name()
	}

	convert, ok := c[typeName]
	if !ok {
		return ErrConvertFailed{
			message:        "type not supported",
			source:         source,
			targetTypeName: typeName,
		}
	}
	err := convert(source, target)
	if err != nil {
		return ErrConvertFailed{
			message:        err.Error(),
			source:         source,
			targetTypeName: typeName,
		}
	}
	return err
}

func Define[T any](l Provider, key string) T {
	var value T
	raw, ok := l.Get(key)
	if !ok {
		l.NotifyError(key, fmt.Errorf("value %s not found", key))
		return value
	}

	valuePtr := reflect.ValueOf(&value).Elem()
	err := supportedConverters.convert(raw.Val, valuePtr)
	if err != nil {
		l.NotifyError(key, fmt.Errorf("error converting path %s: %w", key, err))
	}
	return value
}
