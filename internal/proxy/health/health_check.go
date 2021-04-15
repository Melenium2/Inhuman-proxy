package health

import (
	"context"
	"net"
	"time"
)

// Check checks is proxy still available otherwise returning error
func Check(address string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := net.Dial("tcp", address)
	if err != nil {
		return ErrProxyUnreachable
	}

	if ctx.Err() != nil {
		return ErrProxyTimout(address)
	}

	return nil
}
