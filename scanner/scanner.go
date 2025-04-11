// scanner/scanner.go
package scanner

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sbelyanin/SSCAN/config"
	"github.com/sbelyanin/SSCAN/metrics"
	"github.com/sirupsen/logrus"
)

type ScanTask struct {
	Path   string
	Period time.Duration
}

type Scanner struct {
	ctx   context.Context
	tasks []ScanTask
	mu    sync.RWMutex
}

func NewScanner(ctx context.Context, scanConfigs []config.ScanConfig) *Scanner {
	tasks := make([]ScanTask, len(scanConfigs))
	for i, c := range scanConfigs {
		tasks[i] = ScanTask{
			Path:   c.Path,
			Period: time.Duration(c.Period) * time.Second,
		}
	}
	return &Scanner{
		ctx:   ctx,
		tasks: tasks,
	}
}

func (s *Scanner) UpdateConfig(newTasks []config.ScanConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tasks := make([]ScanTask, len(newTasks))
	for i, c := range newTasks {
		tasks[i] = ScanTask{
			Path:   c.Path,
			Period: time.Duration(c.Period) * time.Second,
		}
	}
	s.tasks = tasks
}

func (s *Scanner) Start() {
	s.mu.RLock()
	tasks := s.tasks
	s.mu.RUnlock()

	wg := sync.WaitGroup{}
	for _, task := range tasks {
		wg.Add(1)
		go func(task ScanTask) {
			defer wg.Done()
			ticker := time.NewTicker(task.Period)
			defer ticker.Stop()
			for {
				select {
				case <-s.ctx.Done():
					return
				case <-ticker.C:
					s.scanPath(task.Path)
				}
			}
		}(task)
	}
	wg.Wait()
}

func (s *Scanner) scanPath(path string) {
	logrus.Debugf("Scanning path %s", path)
	var files []string
	var err error
	if isFile(path) {
		files = []string{path}
	} else {
		files, err = getCertificateFiles(path)
		if err != nil {
			logrus.Warnf("Failed to list files in %s: %v", path, err)
			return
		}
	}

	for _, file := range files {
		err := s.checkCertificate(file)
		if err != nil {
			logrus.Warnf("Error processing %s: %v", file, err)
		}
	}
}

func (s *Scanner) checkCertificate(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return fmt.Errorf("failed to parse PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	metrics.CertExpiryTime.WithLabelValues(path).Set(float64(cert.NotAfter.Unix()))
	metrics.CertSubject.WithLabelValues(path, cert.Subject.String()).Set(1)
	metrics.CertIssuer.WithLabelValues(path, cert.Issuer.String()).Set(1)
	metrics.CertNotBefore.WithLabelValues(path).Set(float64(cert.NotBefore.Unix()))
	metrics.CertNotAfter.WithLabelValues(path).Set(float64(cert.NotAfter.Unix()))

	return nil
}

func isFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}

func getCertificateFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if isCertificateFile(path) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func isCertificateFile(filename string) bool {
	extensions := []string{".pem", ".crt", ".cer"}
	ext := filepath.Ext(filename)
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}
