package kernel

import (
	"net"
)

type Handler func(conn net.Conn) error

type IProxy interface {
	// SetIP Set listening address
	SetIP(ip string)
	// SetPort Set listening port
	SetPort(port int)
	// SetAuth Set the user and password
	SetAuth(auth map[string]string)
	Start() error
	HandleConnection(conn net.Conn) error
}
