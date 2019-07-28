package provision

import (
	"testing"

	"github.com/metrue/fx/config"
)

func TestStart(t *testing.T) {
	host := config.Host{Host: "127.0.0.1"}
	provisionor := New(host)
	if err := provisionor.Start(); err != nil {
		t.Fatal(err)
	}
}
