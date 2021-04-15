package main

import (
	"github.com/Melenium2/inhuman-reverse-proxy/internal/config"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/proxy"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/proxy/storage"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func main() {
	logger, _ := zap.NewProduction(
		zap.AddCaller(),
	)
	defer logger.Sync() //nolint:errcheck
	log := logger.Sugar()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := storage.Connect(cfg.StorageConfig)
	if err != nil {
		log.Fatal(err)
	}
	store := storage.New(conn, log)
	checker := storage.NewChecker(store, log, cfg.CheckerConfig)
	go checker.Check()

	server := proxy.New(cfg, store, log)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig

	if err := server.Shutdown(); err != nil {
		log.Fatal(err)
	}

	log.Info("proxy server shutdown")
}
