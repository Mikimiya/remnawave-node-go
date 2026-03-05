package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hteppl/remnawave-node-go/internal/api/httputil"
)

func PortGuardMiddleware(expectedPort int) gin.HandlerFunc {
	return func(c *gin.Context) {
		localAddr := c.Request.Context().Value(http.LocalAddrContextKey)
		if localAddr == nil {
			httputil.DestroySocket(c)
			return
		}

		tcpAddr, ok := localAddr.(*net.TCPAddr)
		if !ok {
			httputil.DestroySocket(c)
			return
		}

		if tcpAddr.Port != expectedPort || tcpAddr.IP.String() != "127.0.0.1" {
			httputil.DestroySocket(c)
			return
		}

		c.Next()
	}
}
