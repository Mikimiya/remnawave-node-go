package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/hteppl/remnawave-node-go/internal/api/httputil"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestDestroySocket_AbortsContext(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		httputil.DestroySocket(c)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// httptest.ResponseRecorder doesn't implement Hijacker, so DestroySocket
	// will return early from the hijack check but still call c.Abort() via defer
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDestroySocket_NonHijackableWriter_DoesNotPanic(t *testing.T) {
	router := gin.New()

	handlerReached := false
	nextHandlerReached := false

	router.GET("/test", func(c *gin.Context) {
		handlerReached = true
		httputil.DestroySocket(c)
	}, func(c *gin.Context) {
		nextHandlerReached = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Should not panic
	assert.NotPanics(t, func() {
		router.ServeHTTP(w, req)
	})

	assert.True(t, handlerReached, "first handler should be reached")
	assert.False(t, nextHandlerReached, "next handler should NOT be reached after DestroySocket (c.Abort)")
}
