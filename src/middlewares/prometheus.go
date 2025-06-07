package middlewares

import (
	"strconv"
	"time"
	"twitter/src/metrics"

	"github.com/gin-gonic/gin"
)

func Prometheus(ctx *gin.Context) {
	start := time.Now()
	path := ctx.FullPath()
	method := ctx.Request.Method

	ctx.Next()
	status := ctx.Writer.Status()

	metrics.HttpDuration.WithLabelValues(path,method,strconv.Itoa(status)).
	Observe(float64(time.Since(start) / time.Second))
}