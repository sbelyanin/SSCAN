package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Addr    string `yaml:"addr"`
	TLSCert string `yaml:"tls_cert"`
	TLSKey  string `yaml:"tls_key"`
}

func RunServer(ctx context.Context, config ServerConfig) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	authHandler := authMiddleware(config.AuthHashFile, mux)

	server := &http.Server{
		Addr:         config.Addr,
		Handler:      authHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if config.TLSCert != "" && config.TLSKey != "" {
		logrus.Infof("Starting HTTPS server on %s", config.Addr)
		go func() {
			err := server.ListenAndServeTLS(config.TLSCert, config.TLSKey)
			if err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("TLS failed: %v", err)
			}
		}()
	} else {
		logrus.Infof("Starting HTTP server on %s", config.Addr)
		go func() {
			err := server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("ListenAndServe failed: %v", err)
			}
		}()
	}

	<-ctx.Done()
	logrus.Info("Shutting down HTTP server")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(ctxShutdown)
}

func authMiddleware(authFile string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHash, err := getAuthHash(authFile)
		if err != nil {
			logrus.Errorf("Auth hash error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		providedHash := r.Header.Get("X-Auth-Token")
		if providedHash != authHash {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getAuthHash(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read auth file: %w", err)
	}
	return string(data), nil
}
