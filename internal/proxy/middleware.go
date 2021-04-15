package proxy

import uuid "github.com/satori/go.uuid"

const (
	RequestIDHeader = "X-Inhuman-Request-ID"
	ProxyHeader     = "X-Inhuman-Proxy"
)

func generateRequestID() string {
	return uuid.NewV4().String()
}
