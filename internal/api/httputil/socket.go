package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DestroySocket(c *gin.Context) {
	defer func() {
		_ = recover()
		c.Abort()
	}()

	hijacker, ok := c.Writer.(http.Hijacker)
	if !ok {
		return
	}
	conn, _, err := hijacker.Hijack()
	if err != nil {
		return
	}
	conn.Close()
}
