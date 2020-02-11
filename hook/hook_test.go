package hook

import (
	"testing"
)

func TestHook(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		h := New("before_build", "npm install leftpad", "fixture")
		if err := h.Run("fixture"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("script", func(t *testing.T) {
		h := New("before_build", ".hooks/before_build", "fixture")
		if err := h.Run("fixture"); err != nil {
			t.Fatal(err)
		}
	})
}
