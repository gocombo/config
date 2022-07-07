package val

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
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

func unmarshalJSONToStruct(val []byte, target reflect.Value) error {
	newVal := reflect.New(target.Type())
	err := json.Unmarshal(val, newVal.Interface())
	if err != nil {
		return err
	}
	target.Set(reflect.ValueOf(newVal.Elem().Interface()))
	return nil
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
			strSlice = strings.FieldsFunc(actualSliceVal, func(c rune) bool {
				return c == ','
			})
			for i, val := range strSlice {
				strSlice[i] = strings.Trim(val, " ")
			}
		case []string:
			strSlice = actualSliceVal
		case []interface{}: // By default json.Unmarshal will use []interface{} for string arrays
			strSlice = make([]string, len(actualSliceVal))
		parseActualVal:
			for i, val := range actualSliceVal {
				strVal, ok := val.(string)
				if !ok {
					err = fmt.Errorf("expected []string slice value on index=%v: %v(%[2]T)", i, actualSliceVal)
					break parseActualVal
				}
				strSlice[i] = strVal
			}
		default:
			err = fmt.Errorf("expected []string type")
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
		case int32:
			intVal = int(actualVal)
		case int64:
			intVal = int(actualVal)
		case float32:
			actualVal64 := float64(actualVal)
			if math.Trunc(actualVal64) != actualVal64 {
				err = errors.New("float32 value is not an integer")
			}
			intVal = int(actualVal)
		case float64:
			if math.Trunc(actualVal) != actualVal {
				err = errors.New("float64 value is not an integer")
			}
			intVal = int(actualVal)
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
		case float32:
			actualVal64 := float64(actualVal)
			if math.Trunc(actualVal64) != actualVal64 {
				err = errors.New("float32 value is not an integer")
			}
			intVal = int64(actualVal)
		case float64:
			if math.Trunc(actualVal) != actualVal {
				err = errors.New("float64 value is not an integer")
			}
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
		case int32:
			floatVal = float64(actualVal)
		case int64:
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
		default:
			err = errors.New("unexpected bool type")
		}
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(boolVal))
		return nil
	},
	"json-marshaled": func(val interface{}, target reflect.Value) error {
		var jsonData []byte
		switch actualVal := val.(type) {
		case string:
			jsonData = []byte(actualVal)
		default:
			// We by default convert to json first and then unmarshal to target type
			// otherwise it can be quite challenging to handle all individual cases
			var err error
			jsonData, err = json.Marshal(val)
			if err != nil {
				return err
			}
		}
		return unmarshalJSONToStruct(jsonData, target)
	},
	"Duration": func(val interface{}, target reflect.Value) error {
		var durationVal time.Duration
		var err error
		switch actualVal := val.(type) {
		case time.Duration:
			durationVal = actualVal
		case string:
			durationVal, err = time.ParseDuration(actualVal)
		default:
			err = errors.New("unexpected Duration type")
		}
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(durationVal))
		return nil
	},
}

func (c typeConverter) convert(source interface{}, target reflect.Value) error {
	kind := target.Kind()
	targetTypeName := kind.String()

	switch {
	case kind == reflect.Struct || kind == reflect.Map:
		targetTypeName = "json-marshaled"
	case kind == reflect.Slice:
		targetTypeName = "[]" + target.Type().Elem().Name()
	case targetTypeName != target.Type().Name():
		targetTypeName = target.Type().Name()
	}

	convert, ok := c[targetTypeName]
	if !ok {
		return ErrConvertFailed{
			message:        "type not supported",
			source:         source,
			targetTypeName: targetTypeName,
		}
	}
	err := convert(source, target)
	if err != nil {
		return ErrConvertFailed{
			message:        err.Error(),
			source:         source,
			targetTypeName: targetTypeName,
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
