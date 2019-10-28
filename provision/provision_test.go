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

func TestIsFxAgentRunning(t *testing.T) {
	host := config.Host{Host: "127.0.0.1"}
	provisionor := New(host)
	running := provisionor.IsFxAgentRunning()
	if running {
		t.Fatalf("fx-agent should not be running")
	}

	if err := provisionor.StartFxAgent(); err != nil {
		t.Fatal(err)
	}

	running = provisionor.IsFxAgentRunning()
	if !running {
		t.Fatalf("fx-agent should be running")
	}

	if err := provisionor.StopFxAgent(); err != nil {
		t.Fatal(err)
	}
}
