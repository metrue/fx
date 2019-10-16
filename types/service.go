package types

// ServiceRunOptions a service to start options
type ServiceRunOptions struct {
	Image string
	Port  int64
}

// DefaultHost default host IP
const DefaultHost = "0.0.0.0"

// Service instance of a service
type Service struct {
	ID    string `json:"id"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
	State string `json:"state"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
