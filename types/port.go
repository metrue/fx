package types

// PortBinding defines port binding
// ContainerExposePort the port target container exposes
// @ServiceBindingPort the port binding to the port container expose
type PortBinding struct {
	ServiceBindingPort  int32
	ContainerExposePort int32
}
