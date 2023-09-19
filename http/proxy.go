package http

import (
	"fmt"
	"log"
	"net"

	"github.com/pkg6/zproxy/kernel"
)

// http/https 协议
//https://tools.ietf.org/html/draft-luotonen-web-proxy-tunneling-01

const Name = "http"

//curl --proxy "http://127.0.0.1:1080"  http://ipinfo.io/ip
//curl -x "http://admin:123456@127.0.0.1:1080" http://ipinfo.io/ip

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
	kernel.WriteConnectAddr("http", s.Port, s.Auth)
	return kernel.ListenTCP(address, s.HandleConnection)
}

func (s *Proxy) HandleConnection(conn net.Conn) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("http conn handler with panic : %s ", err)
		}
	}()
	req, err := NewHeadBuf(conn, 4096)
	if err != nil {
		req.Close()
		log.Printf("NewHTTPRequest: err=%v", err)
		return err
	}
	if err = s.auth(req); err != nil {
		req.Close()
		log.Printf("http auth: err=%v", err)
		return err
	}
	targetConn, err := s.targetConn(req)
	if err != nil {
		req.Close()
		if targetConn != nil {
			_ = targetConn.Close()
		}
		log.Printf("http connect: err=%v", err)
		return err
	}
	log.Printf("Forwarding from %s to %s", conn.RemoteAddr().String(), targetConn.RemoteAddr().String())
	kernel.Forward(conn, targetConn)
	return nil
}
