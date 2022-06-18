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

func (l *mockLoader) Load(path string) (Raw, error) {
	return l.rawByPath[path], nil
}

func (l *mockLoader) NotifyError(path string, err error) {
	l.errorsByPath[path] = err
}

func TestValue(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		val1Path := fmt.Sprintf("/path1/%s", gofakeit.Word())
		wantVal1Val := gofakeit.SentenceSimple()
		loader := &mockLoader{
			rawByPath: map[string]Raw{
				val1Path: {Val: wantVal1Val},
				fmt.Sprintf("/path2/%s", gofakeit.Word()): {Val: gofakeit.SentenceSimple()},
				fmt.Sprintf("/path3/%s", gofakeit.Word()): {Val: gofakeit.SentenceSimple()},
			},
			errorsByPath: map[string]error{},
		}
		gotVal1Val := Load[string](loader, val1Path)
		assert.Equal(t, wantVal1Val, gotVal1Val)
	})
}
