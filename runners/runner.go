package runners

// Runner make a image a service
type Runner interface {
	Deploy(name string, image string, port int32, svc interface{}) error
	Destroy(name string, svc interface{}) error
	Update(name string, svc interface{}) error
	GetStatus(name string, svc interface{}) error
}
