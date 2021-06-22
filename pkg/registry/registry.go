package registry

import (
	"google.golang.org/grpc/metadata"
)

type ServiceInfo struct {
	Id       string
	Name     string
	Version  string
	Address  string
	Metadata metadata.MD
}

type Registrar interface {
	Register(si *ServiceInfo) error
	Unregister(si *ServiceInfo) error
	Close()
}
