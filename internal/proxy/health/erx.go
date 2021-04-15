package health

import (
	"errors"
	"fmt"
)

var (
	ErrProxyTimout = func(addr string) error {
		return fmt.Errorf("proxy timeout addr = %s", addr)
	}
	ErrProxyUnreachable = errors.New("proxy unreachable")
)
