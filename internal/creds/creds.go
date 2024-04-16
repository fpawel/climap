package creds

import (
	"github.com/fpawel/errorx"
	"strings"
)

type Creds struct {
	Login, Password string
}

func (x Creds) GetCredentials(int) (string, string, error) {
	return x.Login, x.Password, nil
}

func Parse(s string) (Creds, error) {
	xs := strings.Split(s, ",")
	if len(xs) != 2 {
		return Creds{}, errorx.New("bad credentials string")
	}
	return Creds{Login: xs[0], Password: xs[1]}, nil
}
