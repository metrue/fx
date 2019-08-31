package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
}
