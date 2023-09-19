package http

import (
	"fmt"
	"net"

	"github.com/pkg6/zproxy/kernel"
)

func (s *Proxy) targetConn(req ReqHeadBuf) (targetConn net.Conn, err error) {
	inLocalAddr := req.Conn.LocalAddr().String()
	if s.IsDeadLoop(inLocalAddr, req.Host) {
		return nil, fmt.Errorf("dead loop detected , %s", req.Host)
	}
	targetConn, err = net.Dial("tcp", req.Host)
	if err != nil {
		return nil, err
	}
	if req.IsHTTPS() {
		if err = s.HTTPSReply(req); err != nil {
			return nil, err
		}
	} else {
		if _, err = targetConn.Write(req.HeadBuf); err != nil {
			return nil, err
		}
	}
	return targetConn, nil
}

func (s *Proxy) HTTPSReply(req ReqHeadBuf) error {
	_, err := fmt.Fprint(req.Conn, "HTTP/1.1 200 Connection established\r\n\r\n")
	return err
}

func (s *Proxy) IsDeadLoop(inLocalAddr string, host string) bool {
	inIP, inPort, err := net.SplitHostPort(inLocalAddr)
	if err != nil {
		return false
	}
	outDomain, outPort, err := net.SplitHostPort(host)
	if err != nil {
		return false
	}
	if inPort == outPort {
		var outIPs []net.IP
		outIPs, err = net.LookupIP(outDomain)
		if err == nil {
			for _, ip := range outIPs {
				if ip.String() == inIP {
					return true
				}
			}
		}
		if interfaceIPs, err := kernel.GetAllInterfaceAddr(); err == nil {
			for _, localIP := range interfaceIPs {
				for _, outIP := range outIPs {
					if localIP.Equal(outIP) {
						return true
					}
				}
			}
		}
	}
	return false
}
