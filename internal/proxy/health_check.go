package proxy

import (
	"context"
	"net"
	"time"
)

// HealthCheck checks is proxy still available otherwise returning error
func HealthCheck(address string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := net.Dial("tcp", address)
	if err != nil {
		return ErrProxyUnreachable(address)
	}

	if ctx.Err() != nil {
		return ErrProxyTimout(address)
	}

	return nil
}
