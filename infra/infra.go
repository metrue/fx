package infra

// Infra infrastructure provision interface
type Infra interface {
	Provision() (config []byte, err error)
	HealthCheck() (bool, error)
}
