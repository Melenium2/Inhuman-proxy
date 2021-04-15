package main

import (
	"github.com/Melenium2/inhuman-reverse-proxy/internal/config"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/proxy"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/proxy/storage"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	var logger *zap.Logger
	if cfg.DebugMode {
		logger, _ = zap.NewDevelopment(
			zap.AddCaller(),
		)
	} else {
		logger, _ = zap.NewProduction(
			zap.AddCaller(),
		)
	}

	defer logger.Sync() //nolint:errcheck
	log := logger.Sugar()

	conn, err := storage.Connect(cfg.StorageConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("connect to store")

	store := storage.New(conn, log)
	checker := storage.NewChecker(store, log, cfg.CheckerConfig)
	go checker.Check()

	log.Infof("start checking proxies")

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
