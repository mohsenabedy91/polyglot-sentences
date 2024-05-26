package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/metrics"
	"strconv"
	"time"
)

func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		method := c.Request.Method
		c.Next()
		status := c.Writer.Status()
		metrics.HttpDuration.WithLabelValues(path, method, strconv.Itoa(status)).
			Observe(float64(time.Since(start) / time.Millisecond))
	}
}
