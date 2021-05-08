package tcp

import (
	"fmt"
)

/**
 * @Description: 全局ID
 */
type GID struct {
	ServiceName    string
	ServiceVersion string
	IP             string
	Port           int
}

func NewGid(serviceName, serviceVersion, serviceIp string, port int) *GID {
	return &GID{
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		IP:             serviceIp,
		Port:           port,
	}
}

func (gid *GID) String() string {
	return fmt.Sprintf("%v-%v-%v", gid.ServiceName, gid.ServiceVersion, gid.IP)
}

func (gid *GID) NameAndVersion() string {
	return fmt.Sprintf("%v-%v", gid.ServiceName, gid.ServiceVersion)
}
