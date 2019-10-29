package provision

import (
	"testing"
	"time"

	"github.com/metrue/fx/config"
)

func TestProvisionWorkflow(t *testing.T) {
	host := config.Host{Host: "127.0.0.1"}
	provisionor := New(host)

	_ = provisionor.StopFxAgent()
	// TODO wait too long here to make test pass
	time.Sleep(20 * time.Second)

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

	if err := provisionor.Start(); err != nil {
		t.Fatal(err)
	}

	if err := provisionor.StopFxAgent(); err != nil {
		t.Fatal(err)
	}
}
