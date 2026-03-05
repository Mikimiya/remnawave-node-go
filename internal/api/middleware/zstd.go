package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/compress/zstd"
)

func ZstdMiddleware() gin.HandlerFunc {
	decoder, _ := zstd.NewReader(nil)

	return func(c *gin.Context) {
		if c.GetHeader("Content-Encoding") == "zstd" {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(400)
				return
			}
			decompressed, err := decoder.DecodeAll(body, nil)
			if err != nil {
				c.AbortWithStatus(400)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(decompressed))
			c.Request.Header.Del("Content-Encoding")
			c.Request.ContentLength = int64(len(decompressed))
		}
		c.Next()
	}
}
