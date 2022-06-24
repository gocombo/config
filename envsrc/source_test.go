package envsrc

import (
	"fmt"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gocombo/config"
	"github.com/stretchr/testify/assert"
)

type mockLoadOpts struct {
	sourceLoaders []config.SourceLoader
}

func (m *mockLoadOpts) AddSourceLoader(loader config.SourceLoader) {
	m.sourceLoaders = append(m.sourceLoaders, loader)
}

func Test_EnvSource(t *testing.T) {
	loadFromOpts := func(optsSetters ...SourceOpt) (config.Source, error) {
		mockOpts := &mockLoadOpts{}
		loadOpt := Load(optsSetters...)
		loadOpt(mockOpts)
		if len(mockOpts.sourceLoaders) < 1 {
			return nil, fmt.Errorf("no source loader added to opts")
		}
		return mockOpts.sourceLoaders[0]()
	}
	assertVal := func(
		t *testing.T,
		source config.Source,
		key string,
		wantVal string,
	) {
		got, ok := source.GetValue(key)
		if !assert.True(t, ok, "key %s not found", key) {
			return
		}
		assert.Equal(t, wantVal, got.Val)
	}
	t.Run("read configured values from given env vars", func(t *testing.T) {
		env1 := gofakeit.Generate("TEST_ENV_1_{word}")
		env2 := gofakeit.Generate("TEST_ENV_1_{word}")
		path1 := gofakeit.Generate("test/path-1/{word}")
		path2 := gofakeit.Generate("test/path-1/{word}")
		val1 := gofakeit.SentenceSimple()
		val2 := gofakeit.SentenceSimple()
		os.Setenv(env1, val1)
		os.Setenv(env2, val2)
		defer os.Unsetenv(env1)
		defer os.Unsetenv(env2)
		source, err := loadFromOpts(
			Set(path1).From(env1),
			Set(path2).From(env2),
		)
		if !assert.NoError(t, err) {
			return
		}
		assertVal(t, source, path1, val1)
		assertVal(t, source, path2, val2)
	})
	t.Run("handle empty values", func(t *testing.T) {
		env1 := gofakeit.Generate("TEST_ENV_1_{word}")
		env2 := gofakeit.Generate("TEST_ENV_1_{word}")
		path1 := gofakeit.Generate("test/path-1/{word}")
		path2 := gofakeit.Generate("test/path-1/{word}")
		val1 := gofakeit.SentenceSimple()
		os.Setenv(env1, val1)
		os.Setenv(env2, "")
		defer os.Unsetenv(env1)
		defer os.Unsetenv(env2)
		source, err := loadFromOpts(
			Set(path1).From(env1),
			Set(path2).From(env2),
		)
		if !assert.NoError(t, err) {
			return
		}
		assertVal(t, source, path1, val1)
		assertVal(t, source, path2, "")
	})
	t.Run("handle ignore missing values", func(t *testing.T) {
		env1 := gofakeit.Generate("TEST_ENV_1_{word}")
		env2 := gofakeit.Generate("TEST_ENV_1_{word}")
		path1 := gofakeit.Generate("test/path-1/{word}")
		path2 := gofakeit.Generate("test/path-1/{word}")
		val1 := gofakeit.SentenceSimple()
		os.Setenv(env1, val1)
		defer os.Unsetenv(env1)
		source, err := loadFromOpts(
			Set(path1).From(env1),
			Set(path2).From(env2),
		)
		if !assert.NoError(t, err) {
			return
		}
		assertVal(t, source, path1, val1)
	})
}
