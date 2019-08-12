package packer

import "github.com/metrue/fx/types"

// Packer interface
type Packer interface {
	Pack(serviceName string, fn types.ServiceFunctionSource) (types.Project, error)
}
