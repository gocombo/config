package filesrc

import (
	"fmt"
	"os"
	"path/filepath"
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

func setFileValue(t *testing.T, dir, path, val string) {
	if err := os.WriteFile(filepath.Join(dir, path), []byte(val), 0o644); !assert.NoError(t, err) {
		t.FailNow()
	}
}

func Test_FileSource(t *testing.T) {
	tmpDir := t.TempDir()

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
		assert.Equal(t, wantVal, string(got.Val.([]byte)))
	}
	t.Run("read configured values from given files", func(t *testing.T) {
		filePath1 := gofakeit.Generate("test_env_1_{word}")
		filePath2 := gofakeit.Generate("test_env_2_{word}")
		path1 := gofakeit.Generate("test/path-1/{word}")
		path2 := gofakeit.Generate("test/path-2/{word}")
		val1 := gofakeit.SentenceSimple()
		val2 := gofakeit.SentenceSimple()
		setFileValue(t, tmpDir, filePath1, val1)
		setFileValue(t, tmpDir, filePath2, val2)
		source, err := loadFromOpts(
			Set(path1).From(filepath.Join(tmpDir, filePath1)),
			Set(path2).From(filepath.Join(tmpDir, filePath2)),
		)
		if !assert.NoError(t, err) {
			return
		}
		assertVal(t, source, path1, val1)
		assertVal(t, source, path2, val2)
	})
	t.Run("handle empty values", func(t *testing.T) {
		filePath1 := gofakeit.Generate("test_env_1_{word}")
		filePath2 := gofakeit.Generate("test_env_2_{word}")
		path1 := gofakeit.Generate("test/path-1/{word}")
		path2 := gofakeit.Generate("test/path-2/{word}")
		val1 := gofakeit.SentenceSimple()
		setFileValue(t, tmpDir, filePath1, val1)
		setFileValue(t, tmpDir, filePath2, "")
		source, err := loadFromOpts(
			Set(path1).From(filepath.Join(tmpDir, filePath1)),
			Set(path2).From(filepath.Join(tmpDir, filePath2)),
		)
		if !assert.NoError(t, err) {
			return
		}
		assertVal(t, source, path1, val1)
		assertVal(t, source, path2, "")
	})
	t.Run("fail if missing", func(t *testing.T) {
		filePath1 := gofakeit.Generate("test_env_1_{word}")
		path1 := gofakeit.Generate("test/path-1/{word}")
		_, err := loadFromOpts(
			Set(path1).From(filepath.Join(tmpDir, filePath1)),
		)
		if !assert.Error(t, err) {
			return
		}
		assert.True(t, os.IsNotExist(err))
	})
	t.Run("handle ignore missing values", func(t *testing.T) {
		filePath1 := gofakeit.Generate("test_env_1_{word}")
		filePath2 := gofakeit.Generate("test_env_2_{word}")
		path1 := gofakeit.Generate("test/path-1/{word}")
		path2 := gofakeit.Generate("test/path-2/{word}")
		source, err := loadFromOpts(
			Set(path1).From(filepath.Join(tmpDir, filePath1), IgnoreMissing()),
			Set(path2).From(filepath.Join(tmpDir, filePath2), IgnoreMissing()),
		)
		if !assert.NoError(t, err) {
			return
		}
		assertVal(t, source, path1, "")
		assertVal(t, source, path2, "")
	})
}
