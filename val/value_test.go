package val

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

type mockLoader struct {
	rawByPath    map[string]Raw
	errorsByPath map[string]error
}

func (l *mockLoader) Get(path string) (Raw, bool) {
	raw, ok := l.rawByPath[path]
	if !ok {
		return Raw{}, false
	}
	return raw, true
}

func (l *mockLoader) NotifyError(path string, err error) {
	l.errorsByPath[path] = err
}

type valueTestCaseWant struct {
	val interface{}
	err error
}

type valueTestCase struct {
	name     string
	rawValue Raw
	valueTestCaseWant
	define func(l Provider, key string) interface{}
}

func makeValaueTestCase[T any](
	name string,
	rawValue interface{},
	wantVal interface{},
) valueTestCase {
	return valueTestCase{
		name,
		Raw{Val: rawValue},
		valueTestCaseWant{val: wantVal},
		func(l Provider, key string) interface{} {
			return Define[T](l, key)
		},
	}
}

func makeValaueTestCaseErr[T any](
	name string,
	rawValue interface{},
) valueTestCase {
	return valueTestCase{
		name,
		Raw{Val: rawValue},
		valueTestCaseWant{err: ErrConvertFailed{}},
		func(l Provider, key string) interface{} {
			return Define[T](l, key)
		},
	}
}

func TestValue(t *testing.T) {
	t.Run("types", func(t *testing.T) {
		testCases := []func() valueTestCase{
			func() valueTestCase {
				rawVal := gofakeit.SentenceSimple()
				return makeValaueTestCase[string]("string", rawVal, rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Number(1, 100)
				return makeValaueTestCaseErr[string]("string/not a string", rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Number(10, 1000)
				return makeValaueTestCase[int]("int", rawVal, rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Word()
				return makeValaueTestCaseErr[int]("int/not int", rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Int32()
				return makeValaueTestCase[int]("int/from int32", rawVal, int(rawVal))
			},
			func() valueTestCase {
				rawVal := gofakeit.Int64()
				return makeValaueTestCase[int]("int/from int64", rawVal, int(rawVal))
			},
			func() valueTestCase {
				rawVal := float32(gofakeit.Int32())
				return makeValaueTestCase[int]("int/from float32", rawVal, int(rawVal))
			},
			func() valueTestCase {
				rawVal := gofakeit.Float32Range(100, 200)
				return makeValaueTestCaseErr[int]("int/from float32 fractional", rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Number(100, 200)
				return makeValaueTestCase[int]("int/from float64", float64(rawVal), rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Float64Range(100, 200)
				return makeValaueTestCaseErr[int]("int/from float64 fractional", rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Number(10, 1000)
				return makeValaueTestCase[int]("int/from string", strconv.Itoa(rawVal), rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Int64()
				return makeValaueTestCase[int64]("int64", rawVal, rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Number(10, 1000)
				return makeValaueTestCase[int64]("int64/from int", rawVal, int64(rawVal))
			},
			func() valueTestCase {
				rawVal := gofakeit.Int32()
				return makeValaueTestCase[int64]("int64/from int32", rawVal, int64(rawVal))
			},
			func() valueTestCase {
				rawVal := float32(gofakeit.Int32())
				return makeValaueTestCase[int64]("int64/from float32", rawVal, int64(rawVal))
			},
			func() valueTestCase {
				rawVal := gofakeit.Float32Range(100, 200)
				return makeValaueTestCaseErr[int64]("int64/from float32 fractional", rawVal)
			},
			func() valueTestCase {
				rawVal := float64(gofakeit.Int64())
				return makeValaueTestCase[int64]("int64/from float64", rawVal, int64(rawVal))
			},
			func() valueTestCase {
				rawVal := gofakeit.Float64Range(100, 200)
				return makeValaueTestCaseErr[int64]("int64/from float64 fractional", rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Int64()
				return makeValaueTestCase[int64]("int64/from string", strconv.FormatInt(rawVal, 10), rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Bool()
				return makeValaueTestCase[bool]("bool", rawVal, rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Number(10, 100)
				return makeValaueTestCaseErr[bool]("bool/not bool", rawVal)
			},
			func() valueTestCase {
				rawVal := gofakeit.Bool()
				return makeValaueTestCase[bool]("bool/from string", strconv.FormatBool(rawVal), rawVal)
			},

			// TODO:
			/*
			* - duration from string
			* - duration number should fail
			* - complex struct or slice
			* - string[] from csv string
			 */
		}

		for _, tt := range testCases {
			tt := tt()
			t.Run(tt.name, func(t *testing.T) {
				valPath := fmt.Sprintf("/path1/%s", gofakeit.Word())
				loader := &mockLoader{
					rawByPath: map[string]Raw{
						valPath: tt.rawValue,
					},
					errorsByPath: map[string]error{},
				}
				gotVal := tt.define(loader, valPath)
				gotErr := loader.errorsByPath[valPath]
				if tt.valueTestCaseWant.err != nil {
					assert.ErrorAs(t, gotErr, &tt.valueTestCaseWant.err)
					return
				}
				if !assert.NoError(t, gotErr) {
					return
				}
				assert.Equal(t, tt.valueTestCaseWant.val, gotVal)
			})
		}
	})

	t.Run("Define", func(t *testing.T) {
		rawByPath := map[string]Raw{
			fmt.Sprintf("/seed-path1/%s", gofakeit.Word()): {Val: gofakeit.SentenceSimple()},
			fmt.Sprintf("/seed-path2/%s", gofakeit.Word()): {Val: gofakeit.SentenceSimple()},
		}
		loader := &mockLoader{
			rawByPath:    rawByPath,
			errorsByPath: map[string]error{},
		}
		t.Run("existing value", func(t *testing.T) {
			val1Path := fmt.Sprintf("/path1/%s", gofakeit.Word())
			wantVal1Val := gofakeit.SentenceSimple()
			rawByPath[val1Path] = Raw{Val: wantVal1Val}
			gotVal1Val := Define[string](loader, val1Path)
			assert.Equal(t, wantVal1Val, gotVal1Val)
		})
		t.Run("non existing value", func(t *testing.T) {
			val1Path := fmt.Sprintf("/path1/%s", gofakeit.Word())
			gotVal1Val := Define[string](loader, val1Path)
			assert.Equal(t, "", gotVal1Val)
			assert.Len(t, loader.errorsByPath, 1)
			assert.Equal(t, fmt.Errorf("value %s not found", val1Path), loader.errorsByPath[val1Path])
		})
		t.Run("invalid value", func(t *testing.T) {
			val1Path := fmt.Sprintf("/path1/%s", gofakeit.Word())
			rawByPath[val1Path] = Raw{Val: gofakeit.Number(1, 100)}
			gotVal1Val := Define[string](loader, val1Path)
			assert.Equal(t, "", gotVal1Val)
			assert.ErrorIs(t, loader.errorsByPath[val1Path], ErrConvertFailed{})
		})
	})
}
