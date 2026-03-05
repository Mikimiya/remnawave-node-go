package middleware_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hteppl/remnawave-node-go/internal/api/middleware"
)

func TestZstdMiddleware_NoEncoding_Passthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)

	original := `{"key":"value"}`

	var receivedBody string
	router := gin.New()
	router.Use(middleware.ZstdMiddleware())
	router.POST("/test", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		receivedBody = string(body)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(original)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, original, receivedBody)
}

func TestZstdMiddleware_ValidZstd_Decompresses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	original := `{"key":"value"}`

	encoder, err := zstd.NewWriter(nil)
	require.NoError(t, err)
	compressed := encoder.EncodeAll([]byte(original), nil)

	var receivedBody string
	router := gin.New()
	router.Use(middleware.ZstdMiddleware())
	router.POST("/test", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		receivedBody = string(body)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(compressed))
	req.Header.Set("Content-Encoding", "zstd")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, original, receivedBody)
}

func TestZstdMiddleware_InvalidZstd_Returns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.ZstdMiddleware())
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte("not valid zstd data")))
	req.Header.Set("Content-Encoding", "zstd")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestZstdMiddleware_RemovesContentEncodingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	original := `{"test":"data"}`

	encoder, err := zstd.NewWriter(nil)
	require.NoError(t, err)
	compressed := encoder.EncodeAll([]byte(original), nil)

	var contentEncoding string
	var contentLength int64
	router := gin.New()
	router.Use(middleware.ZstdMiddleware())
	router.POST("/test", func(c *gin.Context) {
		contentEncoding = c.GetHeader("Content-Encoding")
		contentLength = c.Request.ContentLength
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(compressed))
	req.Header.Set("Content-Encoding", "zstd")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, contentEncoding, "Content-Encoding should be removed after decompression")
	assert.Equal(t, int64(len(original)), contentLength, "Content-Length should be updated to decompressed size")
}
