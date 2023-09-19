package socks

import (
	"fmt"
	"log"
	"net"

	"github.com/pkg6/zproxy/kernel"
)

// socks 协议
//https://www.ietf.org/archive/id/draft-ietf-aft-socks-chap-01.txt

const (
	Name          = "socks"
	SOCKS5Version = 0x05
	ReservedField = 0x00
)

//curl --proxy "socks5://127.0.0.1:1080"  http://ipinfo.io/ip
//curl --proxy "socks5://admin:123456@127.0.0.1:1080"  http://ipinfo.io/ip

type Proxy struct {
	IP   string
	Port int
	Auth map[string]string
}

func (s *Proxy) SetIP(ip string) {
	s.IP = ip
}
func (s *Proxy) SetPort(port int) {
	s.Port = port
}
func (s *Proxy) SetAuth(auth map[string]string) {
	s.Auth = auth
}

func (s *Proxy) Start() error {
	address := fmt.Sprintf("%s:%d", s.IP, s.Port)
	kernel.WriteConnectAddr("socks5", s.Port, s.Auth)
	return kernel.ListenTCP(address, s.HandleConnection)
}

func (s *Proxy) HandleConnection(conn net.Conn) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("socks conn handler with panic : %s", err)
		}
	}()
	if err := s.auth(conn); err != nil {
		log.Printf("socket auth: err=%v", err)
		return err
	}
	targetConn, err := s.targetConn(conn)
	if err != nil {
		if targetConn != nil {
			_ = targetConn.Close()
		}
		log.Printf("socket connect: err=%v", err)
		return err
	}
	log.Printf("Forwarding from %s to %s", conn.RemoteAddr().String(), targetConn.RemoteAddr().String())
	kernel.Forward(conn, targetConn)
	return nil
}
