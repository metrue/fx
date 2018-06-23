package utils_test

import (
	"os"
	"testing"

	. "github.com/metrue/fx/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestEnsureFile(t *testing.T) {
	fullPath := "/tmp/a/b/c.json"
	err := EnsureFile(fullPath)
	assert.Nil(t, err)

	_, err = os.Stat(fullPath)
	assert.Nil(t, err)

	os.Remove(fullPath)
}
