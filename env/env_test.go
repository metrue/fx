package env

import "testing"

func TestInit(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatal(err)
	}
}
