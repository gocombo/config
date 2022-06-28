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

func TestValue(t *testing.T) {
	t.Run("types", func(t *testing.T) {
		type want struct {
			val interface{}
			err error
		}

		type testCase struct {
			name     string
			rawValue Raw
			want
			define func(l Provider, key string) interface{}
		}

		testCases := []func() testCase{
			func() testCase {
				wantVal := gofakeit.SentenceSimple()
				return testCase{
					"string",
					Raw{Val: wantVal},
					want{
						val: wantVal,
					},
					func(l Provider, key string) interface{} {
						return Define[string](l, key)
					},
				}
			},
			func() testCase {
				wantVal := gofakeit.Number(1, 100)
				return testCase{
					"string/not a string",
					Raw{Val: wantVal},
					want{
						err: ErrBadType,
					},
					func(l Provider, key string) interface{} {
						return Define[string](l, key)
					},
				}
			},
			func() testCase {
				wantVal := gofakeit.Number(10, 1000)
				return testCase{
					"int",
					Raw{Val: wantVal},
					want{
						val: wantVal,
					},
					func(l Provider, key string) interface{} {
						return Define[int](l, key)
					},
				}
			},
			func() testCase {
				wantVal := gofakeit.Word()
				return testCase{
					"int/not int",
					Raw{Val: wantVal},
					want{
						err: ErrBadType,
					},
					func(l Provider, key string) interface{} {
						return Define[int](l, key)
					},
				}
			},
			func() testCase {
				wantVal := gofakeit.Int32()
				return testCase{
					"int/as int32",
					Raw{Val: wantVal},
					want{
						err: ErrBadType,
					},
					func(l Provider, key string) interface{} {
						return Define[int](l, key)
					},
				}
			},
			func() testCase {
				wantVal := gofakeit.Int64()
				return testCase{
					"int/as int64",
					Raw{Val: wantVal},
					want{
						err: ErrBadType,
					},
					func(l Provider, key string) interface{} {
						return Define[int](l, key)
					},
				}
			},
			func() testCase {
				wantVal := gofakeit.Number(10, 1000)
				return testCase{
					"int64/int",
					Raw{Val: wantVal},
					want{
						err: ErrBadType,
					},
					func(l Provider, key string) interface{} {
						return Define[int64](l, key)
					},
				}
			},
			func() testCase {
				wantVal := gofakeit.Number(10, 1000)
				return testCase{
					"int/from string",
					Raw{Val: strconv.Itoa(wantVal)},
					want{
						err: ErrBadType,
					},
					func(l Provider, key string) interface{} {
						return Define[int64](l, key)
					},
				}
			},

			// TODO:
			/*
			* - int from string
			* - duration from string
			* - duration number should fail
			* - bool
			* - bool from string
			* - bool anything else should fail
			* - complex struct or slice
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
				if tt.want.err != nil {
					assert.ErrorIs(t, gotErr, tt.want.err)
					return
				}
				if !assert.NoError(t, gotErr) {
					return
				}
				assert.Equal(t, tt.want.val, gotVal)
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
			assert.ErrorIs(t, loader.errorsByPath[val1Path], ErrBadType)
		})
	})
}
