package docker

import (
	"github.com/metrue/fx/runners"
)

type Docker struct {
}

func (d *Docker) Deploy(name string, image string, port int32, svc interface{}) error {
	return nil
}

func (d *Docker) Update(name string, svc interface{}) error {
	return nil
}

func (d *Docker) Destroy(name string, svc interface{}) error {
	return nil
}

func (d *Docker) GetStatus(name string, svc interface{}) error {
	return nil
}

var (
	_ runners.Runner = &Docker{}
)
