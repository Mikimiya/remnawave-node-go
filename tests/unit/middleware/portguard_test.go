package middleware_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/hteppl/remnawave-node-go/internal/api/middleware"
)

func requestWithLocalAddr(method, path string, addr net.Addr) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	ctx := context.WithValue(req.Context(), http.LocalAddrContextKey, addr)
	return req.WithContext(ctx)
}

func TestPortGuardMiddleware_CorrectPortAndIP_Allows(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var handlerCalled atomic.Bool
	router := gin.New()
	router.Use(middleware.PortGuardMiddleware(61001))
	router.GET("/test", func(c *gin.Context) {
		handlerCalled.Store(true)
		c.Status(http.StatusOK)
	})

	req := requestWithLocalAddr("GET", "/test", &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 61001,
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.True(t, handlerCalled.Load(), "handler should be called for correct port and IP")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPortGuardMiddleware_WrongPort_Blocks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var handlerCalled atomic.Bool
	router := gin.New()
	router.Use(middleware.PortGuardMiddleware(61001))
	router.GET("/test", func(c *gin.Context) {
		handlerCalled.Store(true)
		c.Status(http.StatusOK)
	})

	req := requestWithLocalAddr("GET", "/test", &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9999,
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.False(t, handlerCalled.Load(), "handler should NOT be called for wrong port")
}

func TestPortGuardMiddleware_WrongIP_Blocks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var handlerCalled atomic.Bool
	router := gin.New()
	router.Use(middleware.PortGuardMiddleware(61001))
	router.GET("/test", func(c *gin.Context) {
		handlerCalled.Store(true)
		c.Status(http.StatusOK)
	})

	req := requestWithLocalAddr("GET", "/test", &net.TCPAddr{
		IP:   net.ParseIP("192.168.1.1"),
		Port: 61001,
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.False(t, handlerCalled.Load(), "handler should NOT be called for wrong IP")
}

func TestPortGuardMiddleware_NilLocalAddr_Blocks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var handlerCalled atomic.Bool
	router := gin.New()
	router.Use(middleware.PortGuardMiddleware(61001))
	router.GET("/test", func(c *gin.Context) {
		handlerCalled.Store(true)
		c.Status(http.StatusOK)
	})

	// httptest.NewRequest does not set LocalAddrContextKey by default
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.False(t, handlerCalled.Load(), "handler should NOT be called when LocalAddr is nil")
}

func TestPortGuardMiddleware_NonTCPAddr_Blocks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var handlerCalled atomic.Bool
	router := gin.New()
	router.Use(middleware.PortGuardMiddleware(61001))
	router.GET("/test", func(c *gin.Context) {
		handlerCalled.Store(true)
		c.Status(http.StatusOK)
	})

	req := requestWithLocalAddr("GET", "/test", &net.UnixAddr{
		Name: "/tmp/test.sock",
		Net:  "unix",
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.False(t, handlerCalled.Load(), "handler should NOT be called for non-TCP addr")
}
