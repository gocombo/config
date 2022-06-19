package envsrc

import (
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gocombo/config"
	"github.com/gocombo/config/val"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func Test_EnvSource(t *testing.T) {
	assertVal := func(
		t *testing.T,
		values []val.Raw,
		key string,
		wantVal string,
	) {
		foundIndex := slices.IndexFunc(values, func(r val.Raw) bool {
			return r.Key == key
		})
		if !assert.NotEqual(t, -1, foundIndex, "%s not found", key) {
			return
		}
		assert.Equal(t, wantVal, values[foundIndex].Val)
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
		source := New(
			Set(path1).From(env1),
			Set(path2).From(env2),
		)
		vals, err := source.ReadValues([]string{path1, path2})
		if !assert.NoError(t, err) {
			return
		}
		assert.Len(t, vals, 2)
		assertVal(t, vals, path1, val1)
		assertVal(t, vals, path2, val2)
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
		source := New(
			Set(path1).From(env1),
			Set(path2).From(env2),
		)
		vals, err := source.ReadValues([]string{path1, path2})
		if !assert.NoError(t, err) {
			return
		}
		assert.Len(t, vals, 2)
		assertVal(t, vals, path1, val1)
		assertVal(t, vals, path2, "")
	})
	t.Run("handle ignore missing values", func(t *testing.T) {
		env1 := gofakeit.Generate("TEST_ENV_1_{word}")
		env2 := gofakeit.Generate("TEST_ENV_1_{word}")
		path1 := gofakeit.Generate("test/path-1/{word}")
		path2 := gofakeit.Generate("test/path-1/{word}")
		val1 := gofakeit.SentenceSimple()
		os.Setenv(env1, val1)
		defer os.Unsetenv(env1)
		source := New(
			Set(path1).From(env1),
			Set(path2).From(env2),
		)
		vals, err := source.ReadValues([]string{path1, path2})
		if !assert.NoError(t, err) {
			return
		}
		assert.Len(t, vals, 1)
		assertVal(t, vals, path1, val1)
	})
	t.Run("no support of updated values", func(t *testing.T) {
		env1 := gofakeit.Generate("TEST_ENV_1_{word}")
		env2 := gofakeit.Generate("TEST_ENV_1_{word}")
		path1 := gofakeit.Generate("test/path-1/{word}")
		path2 := gofakeit.Generate("test/path-1/{word}")
		val1 := gofakeit.SentenceSimple()
		os.Setenv(env1, val1)
		os.Setenv(env2, "")
		defer os.Unsetenv(env1)
		defer os.Unsetenv(env2)
		source := New(
			Set(path1).From(env1),
			Set(path2).From(env2),
		)
		vals, err := source.ReadValues(
			[]string{path1, path2},
			config.ReadValuesChangedSince(gofakeit.Date()),
		)
		if !assert.NoError(t, err) {
			return
		}
		assert.Len(t, vals, 0)
	})
}
