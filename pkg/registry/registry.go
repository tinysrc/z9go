package registry

type ServiceInfo struct {
	Id      string
	Name    string
	Version string
	Address string
}

type Registrar interface {
	Register(si *ServiceInfo) error
	Unregister(si *ServiceInfo) error
}
