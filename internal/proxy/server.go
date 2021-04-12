package proxy

import (
	"github.com/Melenium2/inhuman-reverse-proxy/internal/config"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/proxy/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

// Server load balancing incoming requests and add proxy header to each request.
type Server struct {
	// Application config
	Config config.Config

	// Storage which contains available proxy list
	proxyStore storage.ProxyStorage
	// Logger
	logger *zap.SugaredLogger
	// fiber.App provide load balancing and routing utils
	app *fiber.App
}

// New creates new instance of proxy server
func New(logger *zap.SugaredLogger, cfg config.Config) *Server {
	return &Server{
		Config: cfg,
		logger: logger,
	}
}

// Start listening incoming requests and init load balancing and additional routes
func (s *Server) Start() error {
	if err := s.init(); err != nil {
		return err
	}

	return s.app.Listen(":19123")
}

// Stop the server
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

// init creates new load balancing middleware and add to each request proxy settings
func (s *Server) init() error {
	s.app = fiber.New()
	s.app.Use(requestid.New(
		requestid.Config{
			Header: "Inhuman-Request-ID",
		},
	))

	// TODO
	//		add logic with proxy
	s.app.Use(proxy.Balancer(proxy.Config{
		Servers: s.nodes(),
		ModifyRequest: func(ctx *fiber.Ctx) error {
			// TODO Get some proxy
			ctx.Request().Header.Add("X-Inhuman-Proxy", "123")
			return nil
		},
		ModifyResponse: func(ctx *fiber.Ctx) error {
			// TODO Remove proxy and return response to user
			ctx.Response().Header.Del("X-Inhuman-Proxy")
			return nil
		},
	}))

	_ = s.routes()

	return nil
}

// Create additional routes
func (s *Server) routes() error {
	return nil
}

// Convert url.URL to []string which represents servers for load balancing
func (s *Server) nodes() []string {
	nodes := make([]string, len(s.Config.Servers))
	for i, server := range s.Config.Servers {
		nodes[i] = server.String()
	}

	return nodes
}
