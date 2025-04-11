// main.go
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sbelyanin/SSCAN/config"
	"github.com/sbelyanin/SSCAN/logger"
	"github.com/sbelyanin/SSCAN/metrics"
	"github.com/sbelyanin/SSCAN/scanner"
	"github.com/sbelyanin/SSCAN/server"

	"github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	// Инициализация конфига
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логгера
	logger.InitLogger(cfg.Logger)
	logrus.Infof("Starting cert-exporter v1.0.0")

	// Инициализация метрик
	metrics.InitMetrics()

	// Создание контекста для управления процессами
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск HTTP сервера
	go func() {
		err := server.RunServer(ctx, cfg.Server)
		if err != nil {
			logrus.Fatalf("Server failed: %v", err)
		}
	}()

	// Запуск сканера
	scanner := scanner.NewScanner(ctx, cfg.Scans)
	scanner.Start()

	// Запуск перечитывания конфига
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.ConfigPeriod) * time.Second)
		for {
			select {
			case <-ticker.C:
				newConfig, err := config.LoadConfig(*configPath)
				if err != nil {
					logrus.Errorf("Failed to reload config: %v", err)
					continue
				}
				scanner.UpdateConfig(newConfig.Scans)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	// Ожидание прерывания
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	logrus.Info("Stopping exporter...")
	cancel()
}
