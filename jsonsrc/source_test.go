package jsonsrc

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gocombo/config/val"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

type closableBuffer bytes.Buffer

func (b *closableBuffer) Read(p []byte) (n int, err error) {
	return (*bytes.Buffer)(b).Read(p)
}

func (b *closableBuffer) Close() error {
	return nil
}

func TestJsonSource(t *testing.T) {
	type mockNested struct {
		StrVal1 string `json:"str_val_1"`
		StrVal2 string `json:"str_val_2"`
	}
	type mockSourceValues struct {
		StrVal1 string     `json:"str_val_1"`
		StrVal2 string     `json:"str_val_2"`
		Nested  mockNested `json:"nested"`
	}

	randomMockSourceValues := func() mockSourceValues {
		return mockSourceValues{
			StrVal1: gofakeit.Word(),
			StrVal2: gofakeit.Word(),
			Nested: mockNested{
				StrVal1: gofakeit.Word(),
				StrVal2: gofakeit.Word(),
			},
		}
	}

	withMockValues := func(mockValues mockSourceValues) SourceOpt {
		return func(s *source) {
			s.openFile = func(fileName string) (file io.ReadCloser, err error) {
				vals, err := json.Marshal(mockValues)
				if err != nil {
					return nil, err
				}
				return (*closableBuffer)(bytes.NewBuffer(vals)), nil
			}
		}
	}

	t.Run("ReadValues", func(t *testing.T) {
		t.Run("should return values from json source", func(t *testing.T) {
			wantFileName := gofakeit.Generate("{name}.json")
			mockValues := randomMockSourceValues()
			source := New(wantFileName, withMockValues(mockValues))
			wantKeys := []string{
				"str_val_1",
				"str_val_2",
				"nested/str_val_1",
				"nested/str_val_2",
			}
			values, err := source.ReadValues(wantKeys)
			if !assert.NoError(t, err) {
				return
			}
			assertVal := func(key string, wantVal string) {
				foundIndex := slices.IndexFunc(values, func(r val.Raw) bool {
					return r.Key == key
				})
				if !assert.NotEqual(t, -1, foundIndex, "%s not found", key) {
					return
				}
				assert.Equal(t, wantVal, values[foundIndex].Val)
			}
			assertVal("str_val_1", mockValues.StrVal1)
			assertVal("str_val_2", mockValues.StrVal2)
			assertVal("nested/str_val_1", mockValues.Nested.StrVal1)
			assertVal("nested/str_val_2", mockValues.Nested.StrVal2)
		})
		t.Run("ignore non existing values", func(t *testing.T) {
			mockValues := randomMockSourceValues()
			source := New("test.json", withMockValues(mockValues))
			wantKeys := []string{
				"str_val_2",
				"nested/str_val_3",
			}
			values, err := source.ReadValues(wantKeys)
			if !assert.NoError(t, err) {
				return
			}
			assert.Len(t, values, 1)
			foundIndex := slices.IndexFunc(values, func(r val.Raw) bool {
				return r.Key == "nested/srt_val_3"
			})
			assert.Equal(t, -1, foundIndex)
		})
		t.Run("fail if no such file", func(t *testing.T) {
			source := New(gofakeit.Generate("{name}.json"))
			_, err := source.ReadValues([]string{"str_val_1"})
			if !assert.Error(t, err) {
				return
			}
			assert.ErrorIs(t, err, os.ErrNotExist)
		})
		t.Run("ignore missing file", func(t *testing.T) {
			source := New(gofakeit.Generate("{name}.json"), IgnoreMissingFile())
			vals, err := source.ReadValues([]string{"str_val_1"})
			if assert.NoError(t, err) {
				return
			}
			assert.Len(t, vals, 0)
		})
		t.Run("fail if not a JSON", func(t *testing.T) {
			source := New(gofakeit.Generate("{name}.json"), func(opts *source) {
				opts.openFile = func(fileName string) (file io.ReadCloser, err error) {
					return (*closableBuffer)(bytes.NewBufferString("not a json")), nil
				}
			})
			_, err := source.ReadValues([]string{"str_val_1"})
			if !assert.Error(t, err) {
				return
			}
			jsonErr := &json.SyntaxError{}
			assert.ErrorAs(t, err, &jsonErr)
		})
		t.Run("load from given file", func(t *testing.T) {
			wantFileName := gofakeit.Generate("{name}.json")
			var gotFilePath string
			source := New(wantFileName, func(opts *source) {
				opts.openFile = func(fileName string) (file io.ReadCloser, err error) {
					gotFilePath = fileName
					return (*closableBuffer)(bytes.NewBufferString("{}")), nil
				}
			})
			_, err := source.ReadValues([]string{"str_val_1"})
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, wantFileName, gotFilePath)
		})
		t.Run("load from base dir", func(t *testing.T) {
			wantFileName := gofakeit.Generate("{name}.json")
			wantDir := gofakeit.Generate("/{name}/{name}")
			var gotFilePath string
			source := New(wantFileName, WithBaseDir(wantDir), func(opts *source) {
				opts.openFile = func(fileName string) (file io.ReadCloser, err error) {
					gotFilePath = fileName
					return (*closableBuffer)(bytes.NewBufferString("{}")), nil
				}
			})
			_, err := source.ReadValues([]string{"str_val_1"})
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, path.Join(wantDir, wantFileName), gotFilePath)
		})
	})
}
