package http

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg6/zproxy/kernel"
)

func (s *Proxy) auth(req ReqHeadBuf) error {
	if len(s.Auth) != 0 {
		auth, err := s.readerAuth(req)
		if err != nil {
			s.writeAuthUnauthorizedMessage(req)
			return err
		}
		if err := auth.Check(s.Auth); err != nil {
			s.writeAuthUnauthorizedMessage(req)
			return err
		} else {
			return nil
		}
	}
	return nil
}

func (s *Proxy) writeAuthUnauthorizedMessage(req ReqHeadBuf) {
	_, _ = fmt.Fprint(req.Conn, "HTTP/1.1 401 Unauthorized\r\nWWW-Authenticate: Basic realm=\"\"\r\n\r\nUnauthorized")
}

func (s *Proxy) readerAuth(req ReqHeadBuf) (*kernel.Auth, error) {
	auth, err := req.GetHeader("Proxy-Authorization")
	if err != nil {
		return nil, err
	}
	authByte, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(auth, "Basic ", ""))
	if err != nil {
		return nil, err
	}
	authStrings := strings.Split(strings.Trim(string(authByte), " "), ":")
	return &kernel.Auth{Username: authStrings[0], Password: authStrings[1]}, nil
}
