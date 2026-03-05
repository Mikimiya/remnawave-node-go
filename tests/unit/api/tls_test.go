package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hteppl/remnawave-node-go/internal/api"
	"github.com/hteppl/remnawave-node-go/internal/logger"
)

func TestTLSErrorFilter_SuppressesTLSHandshakeErrors(t *testing.T) {
	log := logger.New(logger.Config{Level: logger.LevelError, Format: logger.FormatJSON})
	filter := api.NewTLSErrorFilter(log)

	n, err := filter.Write([]byte("http: TLS handshake error from 1.2.3.4:1234: EOF"))
	require.NoError(t, err)
	assert.Greater(t, n, 0)
}

func TestTLSErrorFilter_PassesOtherErrors(t *testing.T) {
	log := logger.New(logger.Config{Level: logger.LevelError, Format: logger.FormatJSON})
	filter := api.NewTLSErrorFilter(log)

	msg := "some other server error"
	n, err := filter.Write([]byte(msg))
	require.NoError(t, err)
	assert.Equal(t, len(msg), n)
}

func TestTLSErrorFilter_NilLogger(t *testing.T) {
	filter := api.NewTLSErrorFilter(nil)

	// Should not panic with nil logger
	n, err := filter.Write([]byte("some error message"))
	require.NoError(t, err)
	assert.Greater(t, n, 0)

	// TLS handshake errors should also be handled with nil logger
	n, err = filter.Write([]byte("TLS handshake error"))
	require.NoError(t, err)
	assert.Greater(t, n, 0)
}
