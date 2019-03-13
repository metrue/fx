package types

// ServiceRunOptions a service to start options
type ServiceRunOptions struct {
	Image string
	Port  int64
}

// ServiceFunctionSource source of service's function
type ServiceFunctionSource struct {
	Language string `json:"language"`
	Source   string `json:"source"`
}

// ServiceStatus status of service 0, 1, 2
type ServiceStatus int

const (
	// ServiceStatusInit when service initialize
	ServiceStatusInit ServiceStatus = iota
	// ServiceStatusRunning when service is running
	ServiceStatusRunning
	// ServiceStatusStopped when service is stopped
	ServiceStatusStopped
)

// DefaultHost default host IP
const DefaultHost = "0.0.0.0"

// Service a service
type Service struct {
	Name      string        `json:"name"`
	Image     string        `json:"image"`
	Status    ServiceStatus `json:"status"`
	Instances []Instance    `json:"instances"`
}

// Instance instance of a service
type Instance struct {
	ID    string `json:"id"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
	State string `json:"state"`
}
