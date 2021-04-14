package proxy

import "fmt"

var (
	ErrProxyTimout = func(addr string) error {
		return fmt.Errorf("proxy timeout addr = %s", addr)
	}
	ErrProxyUnreachable = func(addr string) error {
		return fmt.Errorf("proxy with addr %s unreachable", addr)
	}
)
