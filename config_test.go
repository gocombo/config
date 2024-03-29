package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gocombo/config/val"
	"github.com/stretchr/testify/assert"
)

type mockKeyValueSource struct {
	values map[string]val.Raw
}

func (m *mockKeyValueSource) GetValue(key string) (val.Raw, bool) {
	v, ok := m.values[key]
	return v, ok
}

func TestLoad(t *testing.T) {
	type config struct {
		val1 string
		val2 string
		val3 string
	}
	createRandomConfig := func() *config {
		return &config{
			val1: gofakeit.Generate("val1-{word}"),
			val2: gofakeit.Generate("val2-{word}"),
			val3: gofakeit.Generate("val3-{word}"),
		}
	}

	testConfigFactory := func(p val.Provider) *config {
		return &config{
			val1: val.Define[string](p, "val1"),
			val2: val.Define[string](p, "val2"),
			val3: val.Define[string](p, "val3"),
		}
	}

	withMockSource := func(src *mockKeyValueSource, err error) LoadOpt {
		return func(opts LoadOpts) {
			opts.AddSourceLoader(func() (Source, error) {
				return src, err
			})
		}
	}

	t.Run("load and build config", func(t *testing.T) {
		want := createRandomConfig()
		got, err := Load(
			testConfigFactory,
			withMockSource(&mockKeyValueSource{
				values: map[string]val.Raw{
					"val1": {Key: "val1", Val: want.val1},
					"val2": {Key: "val2", Val: want.val2},
				},
			}, nil),
			withMockSource(&mockKeyValueSource{
				values: map[string]val.Raw{
					"val3": {Key: "val3", Val: want.val3},
				},
			}, nil),
		)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, want, got)
	})
	t.Run("resolve values in order (last one wins)", func(t *testing.T) {
		want := createRandomConfig()
		got, err := Load(
			testConfigFactory,
			withMockSource(&mockKeyValueSource{
				values: map[string]val.Raw{
					"val1": {Key: "val1", Val: fmt.Sprintf("source-1-%s", want.val1)},
					"val2": {Key: "val2", Val: fmt.Sprintf("source-1-%s", want.val2)},
					"val3": {Key: "val2", Val: fmt.Sprintf("source-1-%s", want.val3)},
				},
			}, nil),
			withMockSource(&mockKeyValueSource{
				values: map[string]val.Raw{
					"val1": {Key: "val1", Val: want.val1},
					"val2": {Key: "val2", Val: want.val2},
					"val3": {Key: "val3", Val: want.val3},
				},
			}, nil),
		)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, want, got)
	})
	t.Run("fail if source failed to load", func(t *testing.T) {
		wantErr := errors.New(gofakeit.Sentence(3))
		_, gotErr := Load(
			testConfigFactory,
			withMockSource(&mockKeyValueSource{
				values: map[string]val.Raw{},
			}, nil),
			withMockSource(&mockKeyValueSource{
				values: map[string]val.Raw{},
			}, wantErr),
		)
		if !assert.Error(t, gotErr) {
			return
		}
		assert.Equal(t, wantErr, gotErr)
	})
	t.Run("fail if notified errors", func(t *testing.T) {
		_, gotErr := Load(
			testConfigFactory,
			withMockSource(&mockKeyValueSource{
				values: map[string]val.Raw{},
			}, nil),
		)
		if !assert.Error(t, gotErr) {
			return
		}
		assert.EqualError(t, gotErr, "failed building config: value val1 not found; value val2 not found; value val3 not found")
	})
	t.Run("fail if no sources", func(t *testing.T) {
		_, err := Load(
			testConfigFactory,
		)
		if assert.Error(t, err) {
			return
		}
		assert.Errorf(t, err, "no sources")
	})
}
