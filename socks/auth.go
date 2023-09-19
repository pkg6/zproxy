package socks

import (
	"fmt"
	"io"

	"github.com/pkg6/zproxy/kernel"
)

const (
	MethodNoAuth       byte = 0x00
	MethodGSSAPI       byte = 0x01
	MethodPassword     byte = 0x02
	MethodNoAcceptable byte = 0xff

	PasswordMethodVersion = 0x01
	PasswordAuthSuccess   = 0x00
	PasswordAuthFailure   = 0x01
)

type AuthMethodReader struct {
	Version  byte
	NMethods byte
	Methods  []byte
}

func (s *Proxy) auth(conn io.ReadWriter) error {
	// Read version, nMethods
	buf := make([]byte, 2)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return err
	}
	// Validate version
	if buf[0] != SOCKS5Version {
		return fmt.Errorf("protocol version not supported")
	}
	// Read methods
	nMethods := buf[1]
	buf = make([]byte, nMethods)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return err
	}
	authReader := &AuthMethodReader{Version: SOCKS5Version, NMethods: nMethods, Methods: buf}
	//By default, no account password is required
	autoMethod := MethodNoAuth
	if len(s.Auth) != 0 {
		//You need to configure the account password
		autoMethod = MethodPassword
	}
	var acceptable bool
	for _, method := range authReader.Methods {
		if method == autoMethod {
			acceptable = true
		}
	}
	if !acceptable {
		_ = s.writerAuthMessage(conn, MethodNoAcceptable)
		return fmt.Errorf("method error reporting is not supported")
	}
	if err = s.writerAuthMessage(conn, autoMethod); err != nil {
		return err
	}
	if autoMethod == MethodPassword {
		auth, err := s.readerAuth(conn)
		if err != nil {
			return err
		}
		if err := auth.Check(s.Auth); err != nil {
			_ = s.writerAuthPasswordMessage(conn, PasswordAuthFailure)
			return err
		} else {
			return s.writerAuthPasswordMessage(conn, PasswordAuthSuccess)
		}
	}
	return nil
}

func (s *Proxy) readerAuth(conn io.Reader) (*kernel.Auth, error) {
	// Read version and username length
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	version, usernameLen := buf[0], buf[1]
	if version != PasswordMethodVersion {
		return nil, ErrMethodVersionNotSupported
	}
	// Read username, password length
	buf = make([]byte, usernameLen+1)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	username, passwordLen := string(buf[:len(buf)-1]), buf[len(buf)-1]
	// Read password
	if len(buf) < int(passwordLen) {
		buf = make([]byte, passwordLen)
	}
	if _, err := io.ReadFull(conn, buf[:passwordLen]); err != nil {
		return nil, err
	}
	return &kernel.Auth{
		Username: username,
		Password: string(buf[:passwordLen]),
	}, nil
}

func (s *Proxy) writerAuthMessage(conn io.Writer, method byte) error {
	_, err := conn.Write([]byte{SOCKS5Version, method})
	return err
}

func (s *Proxy) writerAuthPasswordMessage(conn io.Writer, status byte) error {
	_, err := conn.Write([]byte{PasswordMethodVersion, status})
	return err
}
