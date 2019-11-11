package provision

import (
	"testing"
	"time"
)

func TestProvisionWorkflow(t *testing.T) {
	provisionor := NewWithHost("127.0.0.1", "", "")

	_ = provisionor.StopFxAgent()
	// TODO wait too long here to make test pass
	time.Sleep(40 * time.Second)

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
