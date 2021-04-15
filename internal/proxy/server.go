package proxy

import (
	"fmt"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/config"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/proxy/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
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
func New(cfg config.Config, store storage.ProxyStorage, logger *zap.SugaredLogger) *Server {
	return &Server{
		Config:     cfg,
		proxyStore: store,
		logger:     logger,
	}
}

// Start listening incoming requests and init load balancing and additional routes
func (s *Server) Start() error {
	if err := s.init(); err != nil {
		return err
	}

	return s.app.Listen(fmt.Sprintf(":%d", s.Config.Port))
}

// Shutdown the server
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

// init creates new load balancing middleware and add to each request proxy settings
func (s *Server) init() error {
	s.app = fiber.New()
	s.app.Use(logger.New())

	s.app.Use(
		proxy.Balancer(proxy.Config{
			Servers: s.nodes(),
			ModifyRequest: func(ctx *fiber.Ctx) error {
				address, err := s.proxyStore.GetRandom(ctx.Context())
				if err != nil {
					s.logger.Infof("skip proxy, proxy store return err = %s", err)
				}

				ctx.Request().Header.Add(ProxyHeader, address)
				ctx.Request().Header.Add(RequestIDHeader, generateRequestID())

				return nil
			},
			ModifyResponse: func(ctx *fiber.Ctx) error {
				ctx.Response().Header.Del(ProxyHeader)

				return nil
			},
		}))

	_ = s.routes()

	return nil
}

// Create additional routes
func (s *Server) routes() error {
	s.app.Get("/new/proxy", newProxy(s.proxyStore))

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
