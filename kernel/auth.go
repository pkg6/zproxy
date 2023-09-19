package kernel

import (
	"fmt"
	"log"
)

var (
	ErrPasswordCheckerNotSet = fmt.Errorf("error password checker not set")
	ErrPasswordAuthFailure   = fmt.Errorf("error authenticating username/password")
)

type IAuthCheck interface {
	Check(auth map[string]string) error
}

// Auth Read the incoming account password
type Auth struct {
	Username, Password string
}

// Check 检查账户信息
func (a *Auth) Check(auth map[string]string) error {
	if len(auth) > 0 {
		log.Printf("Received username:%s password %s, verifying operation in progress", a.Username, a.Password)
		if a.Username == "" || a.Password == "" {
			return ErrPasswordCheckerNotSet
		}
		if password, ok := auth[a.Username]; ok && password == a.Password {
			return nil
		}
		return ErrPasswordAuthFailure
	}
	return nil
}
