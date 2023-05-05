package jsonsrc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gocombo/config"
	"github.com/stretchr/testify/assert"
)

type closableBuffer bytes.Buffer

func (b *closableBuffer) Read(p []byte) (n int, err error) {
	return (*bytes.Buffer)(b).Read(p)
}

func (b *closableBuffer) Close() error {
	return nil
}

type mockLoadOpts struct {
	sourceLoaders []config.SourceLoader
}

func (m *mockLoadOpts) AddSourceLoader(loader config.SourceLoader) {
	m.sourceLoaders = append(m.sourceLoaders, loader)
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

	withMockValues := func(mockValues mockSourceValues) LoadOpt {
		return func(s *loadOpts) {
			s.openFile = func(fileName string) (file io.ReadCloser, err error) {
				vals, err := json.Marshal(mockValues)
				if err != nil {
					return nil, err
				}
				return (*closableBuffer)(bytes.NewBuffer(vals)), nil
			}
		}
	}

	t.Run("load", func(t *testing.T) {
		loadFromOpts := func(fileName string, opts ...LoadOpt) (config.Source, error) {
			mockOpts := &mockLoadOpts{}
			loadOpt := Load(fileName, opts...)
			loadOpt(mockOpts)
			if len(mockOpts.sourceLoaders) < 1 {
				return nil, fmt.Errorf("no source loader added to opts")
			}
			return mockOpts.sourceLoaders[0]()
		}
		t.Run("fail if no such file", func(t *testing.T) {
			_, err := loadFromOpts(gofakeit.Generate("{name}.json"))
			assert.ErrorIs(t, err, os.ErrNotExist)
		})
		t.Run("optionally not fail if no such file", func(t *testing.T) {
			source, err := loadFromOpts(gofakeit.Generate("{name}.json"), IgnoreMissingFile())
			if !assert.NoError(t, err) {
				return
			}
			assert.NotNil(t, source)
		})
		t.Run("fail if not a JSON", func(t *testing.T) {
			_, err := loadFromOpts(
				gofakeit.Generate("{name}.json"),
				func(opts *loadOpts) {
					opts.openFile = func(fileName string) (file io.ReadCloser, err error) {
						return (*closableBuffer)(bytes.NewBufferString("not a json")), nil
					}
				},
			)
			if !assert.Error(t, err) {
				return
			}
			jsonErr := &json.SyntaxError{}
			assert.ErrorAs(t, err, &jsonErr)
		})
		t.Run("load from given file", func(t *testing.T) {
			wantFileName := gofakeit.Generate("{name}.json")
			var gotFilePath string
			_, err := loadFromOpts(
				wantFileName,
				func(opts *loadOpts) {
					opts.openFile = func(fileName string) (file io.ReadCloser, err error) {
						gotFilePath = fileName
						return (*closableBuffer)(bytes.NewBufferString("{}")), nil
					}
				},
			)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, wantFileName, gotFilePath)
		})
		t.Run("load from base dir", func(t *testing.T) {
			wantFileName := gofakeit.Generate("{name}.json")
			wantDir := gofakeit.Generate("/{name}/{name}")
			var gotFilePath string
			_, err := loadFromOpts(
				wantFileName,
				WithBaseDir(wantDir),
				func(opts *loadOpts) {
					opts.openFile = func(fileName string) (file io.ReadCloser, err error) {
						gotFilePath = fileName
						return (*closableBuffer)(bytes.NewBufferString("{}")), nil
					}
				},
			)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, path.Join(wantDir, wantFileName), gotFilePath)
		})
	})

	t.Run("GetValue", func(t *testing.T) {
		t.Run("should return values from json source", func(t *testing.T) {
			wantFileName := gofakeit.Generate("{name}.json")
			mockValues := randomMockSourceValues()
			source, err := load(wantFileName, withMockValues(mockValues))
			if !assert.NoError(t, err) {
				return
			}
			assertVal := func(key string, wantVal string) {
				gotVal, ok := source.GetValue(key)
				if !assert.True(t, ok, "Value %s not found", key) {
					return
				}
				assert.Equal(t, wantVal, gotVal.Val)
			}
			assertVal("str_val_1", mockValues.StrVal1)
			assertVal("str_val_2", mockValues.StrVal2)
			assertVal("nested/str_val_1", mockValues.Nested.StrVal1)
			assertVal("nested/str_val_2", mockValues.Nested.StrVal2)
		})
		t.Run("handle non existing data", func(t *testing.T) {
			source, err := load("test.json", IgnoreMissingFile())
			if !assert.NoError(t, err) {
				return
			}
			_, ok := source.GetValue("not/existing/key")
			if !assert.False(t, ok, "Value not/existing/key found") {
				return
			}
		})
	})
}
