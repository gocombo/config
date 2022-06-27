package val

import (
	"fmt"
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
	t.Run("string value", func(t *testing.T) {
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
			assert.Equal(t, fmt.Errorf("value not a string: %s", val1Path), loader.errorsByPath[val1Path])
		})
	})
}
