package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/hteppl/remnawave-node-go/internal/logger"
)

func (s *Server) buildTLSConfig() (*tls.Config, error) {
	cert, err := tls.X509KeyPair(
		[]byte(s.config.Payload.NodeCertPEM),
		[]byte(s.config.Payload.NodeKeyPEM),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM([]byte(s.config.Payload.CACertPEM)) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

type tlsErrorFilter struct {
	logger *logger.Logger
}

func (f *tlsErrorFilter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if strings.Contains(msg, "TLS handshake error") {
		return len(p), nil
	}
	if f.logger != nil {
		f.logger.Error(msg)
	}
	return len(p), nil
}
