package kernel

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"
	"strings"
)

// WriteConnectAddr 显示连接方式日志
func WriteConnectAddr(proxy string, port int, auths map[string]string) {
	ips, _ := GetAllInterfaceAddr()
	var addr []string
	for _, ip := range ips {
		if len(auths) == 0 {
			addr = append(addr, fmt.Sprintf("%s://%s:%d", proxy, ip.String(), port))
		} else {
			for user, password := range auths {
				addr = append(addr, fmt.Sprintf("%s://%s:%s@%s:%d", proxy, user, password, ip.String(), port))
			}
		}
	}
	log.Printf("%s connection address %s ", proxy, strings.Join(addr, " OR "))
}

// GetAllInterfaceAddr  获取本地ip地址
func GetAllInterfaceAddr() ([]net.IP, error) {
	faces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var addresses []net.IP
	for _, iface := range faces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		adds, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range adds {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ip = ip.To4()
			if ip == nil {
				// not an ipv4 address
				continue
			}
			addresses = append(addresses, ip)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no address Found")
	}
	return addresses, nil
}

// ListenTls 监听ssl服务
func ListenTls(ip string, port int, certBytes, keyBytes []byte) (ln *net.Listener, err error) {
	var cert tls.Certificate
	cert, err = tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		err = errors.New("failed to parse root certificate")
	}
	config := &tls.Config{
		ClientCAs:    clientCertPool,
		ServerName:   "proxy",
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	_ln, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err == nil {
		ln = &_ln
	}
	return
}

// ListenTCP 监听tcp服务
func ListenTCP(address string, handleConnection func(conn net.Conn) error) error {
	var err error
	var listener net.Listener
	listener, err = net.Listen("tcp", address)
	if err == nil {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("ListenTCP , err=%s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			for {
				var conn net.Conn
				conn, err = listener.Accept()
				if err == nil {
					go func() {
						defer func() {
							if e := recover(); e != nil {
								log.Printf("connection handler,err=%s , \ntrace:%s", e, string(debug.Stack()))
							}
						}()
						if err = handleConnection(conn); err != nil {
							log.Printf("handle connection failure from %s: %s", conn.RemoteAddr(), err)
						}
					}()
				} else {
					log.Printf("accept error , ERR:%s", err)
					break
				}
			}
		}()
	}
	return err
}

// Forward 转发
func Forward(conn, target net.Conn) {
	if conn == nil || target == nil {
		return
	}
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	go forward(conn, target)
	go forward(target, conn)
}
