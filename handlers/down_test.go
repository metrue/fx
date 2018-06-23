package handlers_test

import (
	"testing"

	. "github.com/metrue/fx/handlers"
	"github.com/stretchr/testify/assert"
)

func TestDown(t *testing.T) {
	containerId := "hello-container-id-not-exist"
	image := "world-image-name=not-exit"
	_, err := Down(containerId, image)
	assert.Equal(t, RemoveContainerError, err)
}
