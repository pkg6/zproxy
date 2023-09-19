package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	IPv4Length = net.IPv4len
	IPv6Length = net.IPv6len

	CommandConnect byte = 0x01
	CommandBind    byte = 0x02
	CommandUDP     byte = 0x03

	AddressTypeIPv4   byte = 0x01
	AddressTypeDomain byte = 0x03
	AddressTypeIPv6   byte = 0x04
)
const (
	ReplySuccess byte = iota
	ReplyServerFailure
	ReplyConnectionNotAllowed
	ReplyNetworkUnreachable
	ReplyHostUnreachable
	ReplyConnectionRefused
	ReplyTTLExpired
	ReplyCommandNotSupported
	ReplyAddressTypeNotSupported
)

var (
	ErrVersionNotSupported       = fmt.Errorf("protocol version not supported")
	ErrMethodVersionNotSupported = fmt.Errorf("sub-negotiation method version not supported")
	ErrCommandNotSupported       = fmt.Errorf("requst command not supported")
	ErrInvalidReservedField      = fmt.Errorf("invalid reserved field")
	ErrAddressTypeNotSupported   = fmt.Errorf("address type not supported")
)

type TargetReader struct {
	Command  byte
	AddrType byte
	Address  string
	Port     uint16
}

func (s *Proxy) targetConn(conn net.Conn) (net.Conn, error) {
	message, err := s.targetReader(conn)
	if err != nil {
		return nil, err
	}
	// Check if the command is supported
	if message.Command != CommandConnect {
		return nil, s.writeTargetFailureMessage(conn, ReplyCommandNotSupported)
	}
	// Check if the address type is supported
	if message.AddrType == IPv6Length {
		return nil, s.writeTargetFailureMessage(conn, ReplyAddressTypeNotSupported)
	}
	// 请求访问目标TCP服务
	address := fmt.Sprintf("%s:%d", message.Address, message.Port)
	targetConn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, s.writeTargetFailureMessage(conn, ReplyConnectionRefused)
	}
	// Send success reply
	addrValue := targetConn.LocalAddr()
	addr := addrValue.(*net.TCPAddr)
	return targetConn, s.writeTargetSuccessMessage(conn, addr.IP, uint16(addr.Port))
}
func (s *Proxy) writeTargetFailureMessage(conn io.Writer, replyType byte) error {
	_, err := conn.Write([]byte{SOCKS5Version, replyType, ReservedField, AddressTypeIPv4, 0, 0, 0, 0, 0, 0})
	return err
}
func (s *Proxy) writeTargetSuccessMessage(conn io.Writer, ip net.IP, port uint16) error {
	addressType := AddressTypeIPv4
	if len(ip) == IPv6Length {
		addressType = AddressTypeIPv6
	}
	// Write version, reply success, reserved, address type
	_, err := conn.Write([]byte{SOCKS5Version, ReplySuccess, ReservedField, addressType})
	if err != nil {
		return err
	}
	// Write bind IP(IPv4/IPv6)
	if _, err := conn.Write(ip); err != nil {
		return err
	}
	// Write bind port
	buf := make([]byte, 2)
	buf[0] = byte(port >> 8)
	buf[1] = byte(port - uint16(buf[0])<<8)
	_, err = conn.Write(buf)
	return err
}

func (s *Proxy) targetReader(conn io.Reader) (*TargetReader, error) {
	// Read version, command, reserved, address type
	buf := make([]byte, 4)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	version, command, reserved, addrType := buf[0], buf[1], buf[2], buf[3]
	// Check if the fields are valid
	if version != SOCKS5Version {
		return nil, ErrVersionNotSupported
	}
	if command != CommandConnect && command != CommandBind && command != CommandUDP {
		return nil, ErrCommandNotSupported
	}
	if reserved != ReservedField {
		return nil, ErrInvalidReservedField
	}
	if addrType != AddressTypeIPv4 && addrType != AddressTypeIPv6 && addrType != AddressTypeDomain {
		return nil, ErrAddressTypeNotSupported
	}
	var err error
	// Read address and port
	message := TargetReader{Command: command, AddrType: addrType}
	switch addrType {
	case AddressTypeIPv6:
		ipv6 := make(net.IP, IPv6Length)
		_, err = conn.Read(ipv6)
		if err != nil {
			log.Printf("Resolve the IPv6 err %s", err)
		}
		message.Address = ipv6.String()
		log.Printf("Resolve the IPv6 destination address %s", message.Address)
	case AddressTypeIPv4:
		ipv4 := make(net.IP, IPv4Length)
		_, err = conn.Read(ipv4)
		if err != nil {
			log.Printf("Resolve the IPv4 err %s", err)
		}
		message.Address = ipv4.String()
		log.Printf("Resolve the IPv4 destination address %s", message.Address)
	case AddressTypeDomain:
		var domainLen uint8
		err = binary.Read(conn, binary.BigEndian, &domainLen)
		if err != nil {
			log.Printf("Resolve the domain binary.Read err %s", err)
		}
		domain := make([]byte, domainLen)
		_, err = conn.Read(domain)
		if err != nil {
			log.Printf("Resolve the domain err %s", err)
		}
		message.Address = string(domain)
		log.Printf("Resolve the Domain destination address %s", message.Address)
	}
	var port uint16
	binary.Read(conn, binary.BigEndian, &port)
	message.Port = port
	return &message, nil
}
