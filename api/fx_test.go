package api_test

import (
	"bytes"
	"flag"
	"log"
	"os/exec"
	"testing"

	. "github.com/metrue/fx/api"
	"github.com/stretchr/testify/assert"
)

func setup() {
	cmd := exec.Command("./gen.sh")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func teardown() {
	cmd := exec.Command("./clean.sh")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func TestNewClient(t *testing.T) {
	endpoint := ":5050"
	client, conn, err := NewClient(endpoint)
	assert.Equal(t, client, nil)
	assert.Equal(t, conn, nil)
	assert.Equal(t, err, nil)
}

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		setup()
	}

	m.Run()

	if !testing.Short() {
		teardown()
	}
}
