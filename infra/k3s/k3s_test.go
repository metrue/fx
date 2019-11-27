package k3s

import (
	"fmt"
	"os"
	"testing"
)

func TestK3S(t *testing.T) {
	if os.Getenv("K3S_MASTER_IP") == "" ||
		os.Getenv("K3S_MASTER_USER") == "" ||
		os.Getenv("K3S_AGENT_IP") == "" ||
		os.Getenv("K3S_AGENT_USER") == "" {
		t.Skip("skip k3s test since K3S_MASTER_IP, K3S_MASTER_USER and K3S_AGENT_IP, K3S_AGENT_USER not ready")
	}

	master := MasterNode{
		IP:   os.Getenv("K3S_MASTER_IP"),
		User: os.Getenv("K3S_MASTER_USER"),
	}
	agents := []AgentNode{
		AgentNode{
			IP:   os.Getenv("K3S_AGENT_IP"),
			User: os.Getenv("K3S_AGENT_USER"),
		},
	}
	k3s := New(master, agents)
	if err := k3s.SetupMaster(); err != nil {
		t.Fatal(err)
	}

	kubeconfig, err := k3s.GetKubeConfig()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(kubeconfig))

	if _, err := k3s.getToken(); err != nil {
		t.Fatal(err)
	}

	if err := k3s.SetupAgent(); err != nil {
		t.Fatal(err)
	}
}
