package spinner

import (
	"fmt"
	"testing"
	"time"
)

func TestSpinner(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		Start("task 2")
		time.Sleep(1 * time.Second)
		Stop("task 2", fmt.Errorf("error happened"))
	})
	t.Run("success", func(t *testing.T) {
		Start("task 1")
		time.Sleep(1 * time.Second)
		Stop("task 1", nil)
	})
}
