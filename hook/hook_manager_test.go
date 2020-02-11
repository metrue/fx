package hook

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHookManager(t *testing.T) {
	t.Run("descovery in default hookdir .hooks", func(t *testing.T) {
		hooks, err := descovery("")
		if err != nil {
			t.Fatal(err)
		}

		if len(hooks) != 1 {
			t.Fatalf("should have one hook, but got %d", len(hooks))
		}

		if hooks[0].Name() != HookNameBeforeBuild {
			t.Fatalf("should be before_build hook, but got %s", hooks[0].Name())
		}
	})

	t.Run("descovery in empty hookdir", func(t *testing.T) {
		hooks, err := descovery(filepath.Join(os.TempDir(), ".hooks"))
		if err != nil {
			t.Fatal(err)
		}
		if len(hooks) != 0 {
			t.Fatalf("should get 0 hooks, but got %d", len(hooks))
		}
	})

	t.Run("run before_build hook", func(t *testing.T) {
		if err := RunBeforeBuildHook("fixture"); err != nil {
			t.Fatal(err)
		}
	})
}
